import ast

print(ast.dump(ast.parse(open('out.py').read(), filename='out.py', type_comments=True),
    indent=4))
