package main

type Editor interface {
	Insert(position int, text string) Editor
	Delete(offset, length uint) Editor
	Undo() Editor
	Redo() Editor
	String() string
}

type piece struct {
	origin bool
	offset int
	length int
}

type PieceTable struct {
	originBuffer string
	addBuffer    string
	table        []piece
}

type DefaultEditor struct {
	PieceTable
}

func NewEditor(text string) Editor {
	editor := DefaultEditor{
		PieceTable: PieceTable{
			originBuffer: text,
			addBuffer:    "",
			table:        make([]piece, 0),
		},
	}

	editor.table = append(editor.table, piece{
		origin: true,
		offset: 0,
		length: len(text),
	})

	return &editor
}

func (editor *DefaultEditor) length() int {
	length := 0
	for _, v := range editor.table {
		length += v.length
	}
	return length
}

func (editor *DefaultEditor) Insert(position int, text string) Editor {
	editor.addBuffer += text
	newPiece := piece{
		origin: false,
		offset: len(editor.addBuffer) - len(text),
		length: len(text),
	}

	if position == 0 {
		editor.table = append([]piece{newPiece}, editor.table...)
	} else if position >= editor.length() {
		editor.table = append(editor.table, newPiece)
	} else {
		curLength := 0
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
	return editor
}

func (editor *DefaultEditor) Undo() Editor {
	return nil
}

func (editor *DefaultEditor) Redo() Editor {
	return nil
}

func (editor *DefaultEditor) String() string {
	result := ""
	for _, v := range editor.table {
		off := v.offset
		length := v.length
		if v.origin {
			result += editor.originBuffer[off : off+length]
		} else {
			result += editor.addBuffer[off : off+length]
		}
	}
	return result
}
