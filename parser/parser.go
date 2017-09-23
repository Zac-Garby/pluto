package parser

import (
	"fmt"

	"github.com/Zac-Garby/pluto/ast"
	"github.com/Zac-Garby/pluto/lexer"
	"github.com/Zac-Garby/pluto/token"
)

type prefixParser func() ast.Expression
type infixParser func(ast.Expression) ast.Expression

// Parser parses a string into an
// abstract syntax tree
type Parser struct {
	Errors []Error

	lex       func() token.Token
	text      string
	cur, peek token.Token
	prefixes  map[token.Type]prefixParser
	infixes   map[token.Type]infixParser
	argTokens []token.Type
}

// New returns a new parser for the
// given string
func New(text, file string) *Parser {
	p := &Parser{
		lex:    lexer.Lexer(text, file),
		text:   text,
		Errors: []Error{},
	}

	p.prefixes = map[token.Type]prefixParser{
		token.ID:         p.parseID,
		token.Number:     p.parseNum,
		token.True:       p.parseBool,
		token.False:      p.parseBool,
		token.Null:       p.parseNull,
		token.LeftSquare: p.parseArrayOrMap,
		token.String:     p.parseString,
		token.Char:       p.parseChar,
		token.LessThan:   p.parseEmission,
		token.Param:      p.parseParam,

		token.Minus: p.parsePrefix,
		token.Plus:  p.parsePrefix,
		token.Bang:  p.parsePrefix,

		token.LeftParen: p.parseGroupedExpression,
		token.If:        p.parseIfExpression,
		token.BackSlash: p.parseFunctionCall,
		token.LeftBrace: p.parseBlockLiteral,
	}

	p.infixes = map[token.Type]infixParser{
		token.Plus:               p.parseInfix,
		token.Minus:              p.parseInfix,
		token.Star:               p.parseInfix,
		token.Slash:              p.parseInfix,
		token.Equal:              p.parseInfix,
		token.NotEqual:           p.parseInfix,
		token.LessThan:           p.parseInfix,
		token.GreaterThan:        p.parseInfix,
		token.Or:                 p.parseInfix,
		token.And:                p.parseInfix,
		token.BitOr:              p.parseInfix,
		token.BitAnd:             p.parseInfix,
		token.Exp:                p.parseInfix,
		token.FloorDiv:           p.parseInfix,
		token.Mod:                p.parseInfix,
		token.LessThanEq:         p.parseInfix,
		token.GreaterThanEq:      p.parseInfix,
		token.QuestionMark:       p.parseInfix,
		token.AndEquals:          p.parseShorthandAssignment,
		token.BitAndEquals:       p.parseShorthandAssignment,
		token.BitOrEquals:        p.parseShorthandAssignment,
		token.ExpEquals:          p.parseShorthandAssignment,
		token.FloorDivEquals:     p.parseShorthandAssignment,
		token.MinusEquals:        p.parseShorthandAssignment,
		token.ModEquals:          p.parseShorthandAssignment,
		token.OrEquals:           p.parseShorthandAssignment,
		token.PlusEquals:         p.parseShorthandAssignment,
		token.QuestionMarkEquals: p.parseShorthandAssignment,
		token.SlashEquals:        p.parseShorthandAssignment,
		token.StarEquals:         p.parseShorthandAssignment,
		token.Assign:             p.parseAssignExpression,
		token.Dot:                p.parseDotExpression,
		token.Colon:              p.parseQualifiedFunctionCall,
		token.LeftSquare:         p.parseIndexExpression,
	}

	p.argTokens = []token.Type{}

	for k := range p.prefixes {
		if !isBlacklisted(k) {
			p.argTokens = append(p.argTokens, k)
		}
	}

	p.next()
	p.next()

	return p
}

func (p *Parser) peekPrecedence() int {
	if precedence, ok := precedences[p.peek.Type]; ok {
		return precedence
	}

	return lowest
}

func (p *Parser) curPrecedence() int {
	if precedence, ok := precedences[p.cur.Type]; ok {
		return precedence
	}

	return lowest
}

func (p *Parser) next() {
	p.cur = p.peek
	p.peek = p.lex()

	if p.peek.Type == token.Illegal {
		p.Err(
			fmt.Sprintf("illegal token found: `%s`", p.peek.Literal),
			p.peek.Start,
			p.peek.End,
		)
	}
}

// Parse parses an entire program
func (p *Parser) Parse() ast.Program {
	prog := ast.Program{
		Statements: []ast.Statement{},
	}

	for !p.curIs(token.EOF) {
		stmt := p.parseStatement()

		if stmt != nil {
			prog.Statements = append(prog.Statements, stmt)
		}

		p.next()
	}

	return prog
}
