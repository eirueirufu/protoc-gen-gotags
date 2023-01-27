%{
    package tags
%}

%union {
    tok string
    tags Tags
    tag *Tag
}

%start root

%token <tok> ident string_lit colon 

%type <tags> root tags
%type <tag> tag

%%

root: 
    '`' tags '`' {
        yylex.(*Lex).root = $2
    }
    | '`' '`' {
        yylex.(*Lex).root = nil
    }
    | {
        yylex.(*Lex).root = nil
    }

tags: 
    tag {
        $$ = []*Tag{$1}
    }
    | tags tag {
        $$ = append($1, $2)
    }

tag: 
    ident ':' string_lit {
        $$ = &Tag{
            Key: $1,
            Value: $3,
        }
    }

%%