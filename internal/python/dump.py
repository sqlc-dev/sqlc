import ast
import argparse

parser = argparse.ArgumentParser(description='Process some integers.')
parser.add_argument('file', type=argparse.FileType('r', encoding='UTF-8'), help='file to parse')
args = parser.parse_args()

print(ast.dump(ast.parse(args.file.read(), filename='out.py', type_comments=True),
    indent=4))
