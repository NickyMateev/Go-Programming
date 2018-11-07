package main

import "bytes"

type Editor interface {
	// Insert text starting from given position.
	Insert(position uint, text string) Editor

	// Delete length items from offset.
	Delete(offset, length uint) Editor

	// Undo reverts latest change.
	Undo() Editor

	// Redo re-applies latest undone change.
	Redo() Editor

	// String returns complete representation of what a file looks
	// like after all manipulations.
	String() string
}

type piece struct {
	origin bool
	offset uint
	length uint
}

type PieceTable struct {
	originBuffer *bytes.Buffer
	addBuffer    *bytes.Buffer
	table        []piece
}

type DefaultEditor struct {
	PieceTable
	tableHistory [][]piece
}

func NewEditor(text string) Editor {
	editor := DefaultEditor{
		PieceTable: PieceTable{
			originBuffer: bytes.NewBufferString(text),
			addBuffer:    bytes.NewBufferString(""),
			table:        make([]piece, 0),
		},
		tableHistory: make([][]piece, 0),
	}

	editor.table = append(editor.table, piece{
		origin: true,
		offset: 0,
		length: uint(len(text)),
	})

	editor.tableHistory = append(editor.tableHistory, editor.table)
	return &editor
}

func (editor *DefaultEditor) length() uint {
	var length uint
	for _, v := range editor.table {
		length += v.length
	}
	return length
}

// reallocateTable reallocates a new slice of pieces and copies all of the current elements onto it - needed for proper Undo and Redo
func (editor *DefaultEditor) reallocateTable() {
	temp := editor.table
	editor.table = make([]piece, len(temp))
	copy(editor.table, temp)
}

// discardUndoneHistory discards the saved undone history by reallocating a new slice of piece slices
func (editor *DefaultEditor) discardUndoneHistory() {
	if len(editor.tableHistory) < cap(editor.tableHistory) {
		temp := editor.tableHistory
		editor.tableHistory = make([][]piece, len(temp))
		copy(editor.tableHistory, temp)
	}
}

func (editor *DefaultEditor) Insert(position uint, text string) Editor {
	defer func() {
		editor.tableHistory = append(editor.tableHistory, editor.table)
	}()

	editor.reallocateTable()
	editor.discardUndoneHistory()

	editor.addBuffer.WriteString(text)
	newPiece := piece{
		origin: false,
		offset: uint((editor.addBuffer.Len()) - len(text)),
		length: uint(len(text)),
	}

	if position == 0 {
		editor.table = append([]piece{newPiece}, editor.table...)
	} else if position >= editor.length() {
		editor.table = append(editor.table, newPiece)
	} else {
		var curLength uint
		for i, elem := range editor.table {
			if (curLength + elem.length) < position {
				curLength += elem.length
				continue
			}
			editor.table[i].length = position - curLength

			residuePiece := elem
			residuePiece.offset += editor.table[i].length
			residuePiece.length -= editor.table[i].length

			editor.table = append(editor.table, piece{}, piece{})
			copy(editor.table[i+3:], editor.table[i+1:])
			editor.table[i+1] = newPiece
			editor.table[i+2] = residuePiece
			break
		}
	}
	return editor
}

func (editor *DefaultEditor) Delete(offset, length uint) Editor {
	defer func() {
		editor.tableHistory = append(editor.tableHistory, editor.table)
	}()

	totalLength := editor.length()
	if offset >= totalLength {
		return editor
	}
	editor.reallocateTable()
	editor.discardUndoneHistory()

	if offset+length > totalLength {
		length = totalLength - offset
	}
	startPieceIdx, endPieceIdx := editor.affectedPiecesFromDelete(offset, length)
	defer editor.cleanupTable(startPieceIdx, endPieceIdx)

	var firstDeletionIdxInStartPiece = offset - editor.absoluteIdxOfPiece(startPieceIdx)
	var lastDeletionIdxInEndPiece = length - (editor.absoluteIdxOfPiece(endPieceIdx) - offset) - 1

	piecesAffected := (endPieceIdx - startPieceIdx) + 1
	if piecesAffected > 1 {
		editor.table[startPieceIdx].length = firstDeletionIdxInStartPiece
		editor.table[endPieceIdx].offset += lastDeletionIdxInEndPiece + 1
		editor.table[endPieceIdx].length -= lastDeletionIdxInEndPiece + 1
	} else {
		originalPieceLength := editor.table[startPieceIdx].length
		editor.table[startPieceIdx].length = firstDeletionIdxInStartPiece

		residuePieceLength := originalPieceLength - length - editor.table[startPieceIdx].length
		if residuePieceLength > 0 {
			residuePiece := editor.table[startPieceIdx]
			residuePiece.offset += editor.table[startPieceIdx].length + length
			residuePiece.length = residuePieceLength

			editor.table = append(editor.table, piece{})
			copy(editor.table[startPieceIdx+2:], editor.table[startPieceIdx+1:])
			editor.table[startPieceIdx+1] = residuePiece
		}
	}
	return editor
}

// cleanupTable deletes all piece elements from the PieceTable which are within the range defined by startIdx and endIdx as well as the elements at the startIdx and endIdx in the cases where they have zero length
func (editor *DefaultEditor) cleanupTable(startIdx, endIdx uint) {
	indicesToDrop := make([]uint, 0)
	for i := startIdx; i <= endIdx; i++ {
		if (i != startIdx && i != endIdx) || editor.table[i].length < 1 {
			indicesToDrop = append(indicesToDrop, i)
		}
	}
	if len(indicesToDrop) > 0 {
		lastIndexToDrop := indicesToDrop[len(indicesToDrop)-1]
		editor.table = append(editor.table[:indicesToDrop[0]], editor.table[lastIndexToDrop+1:]...)
	}
}

// absoluteIdxOfPiece returns the absolute index (in relation to the whole text in the Editor) of the first symbol of the piece at the specified index in the Piece Table
func (editor *DefaultEditor) absoluteIdxOfPiece(pieceIdx uint) uint {
	var absIdx, i uint
	for ; i < pieceIdx; i++ {
		absIdx += editor.table[i].length
	}
	return absIdx
}

// affectedPiecesFromDelete returns the indices of the first and last affected piece from the Piece Table as a result of the delete operation
func (editor *DefaultEditor) affectedPiecesFromDelete(offset, length uint) (uint, uint) {
	var absoluteIdxOfStartingPiece, startPieceIdx, endPieceIdx uint
	for i, elem := range editor.table {
		if (absoluteIdxOfStartingPiece + elem.length) < offset {
			absoluteIdxOfStartingPiece += elem.length
			continue
		}
		startPieceIdx = uint(i)
		break
	}

	deleteSymbolsCounter := (absoluteIdxOfStartingPiece + editor.table[startPieceIdx].length) - offset
	if deleteSymbolsCounter >= length {
		return startPieceIdx, startPieceIdx
	}

	for i := startPieceIdx + 1; i < uint(len(editor.table)); i++ {
		if deleteSymbolsCounter+editor.table[i].length < length {
			deleteSymbolsCounter += editor.table[i].length
			continue
		}
		endPieceIdx = uint(i)
		break
	}
	return startPieceIdx, endPieceIdx
}

func (editor *DefaultEditor) Undo() Editor {
	historyRecords := len(editor.tableHistory)
	if historyRecords > 1 {
		editor.tableHistory = editor.tableHistory[:historyRecords-1]
		editor.table = editor.tableHistory[historyRecords-2]
	}
	return editor
}

func (editor *DefaultEditor) Redo() Editor {
	historyRecords := len(editor.tableHistory)
	if historyRecords < cap(editor.tableHistory) {
		editor.tableHistory = editor.tableHistory[:historyRecords+1]
		editor.table = editor.tableHistory[historyRecords]
	}
	return editor
}

func (editor *DefaultEditor) String() string {
	origin := editor.originBuffer.String()
	add := editor.addBuffer.String()

	var result bytes.Buffer
	for _, v := range editor.table {
		off := v.offset
		length := v.length
		if v.origin {
			result.WriteString(origin[off : off+length])
		} else {
			result.WriteString(add[off : off+length])
		}
	}
	return result.String()
}
