import sys, os, re, shutil
from glob import glob
from pathlib import Path

def main(argv):
    for file in glob("./parser/*Lexer.g4"):
        fix_lexer(file)
    for file in glob("./parser/*Parser.g4"):
        fix_parser(file)

def fix_lexer(file_path):
    print("Altering " + file_path)
    if not os.path.exists(file_path):
        print(f"Could not find file: {file_path}")
        sys.exit(1)
    parts = os.path.split(file_path)
    file_name = parts[-1]

    shutil.move(file_path, file_path + ".bak")
    input_file = open(file_path + ".bak",'r')
    output_file = open(file_path, 'w')
    for x in input_file:
        if 'this.' in x and '}?' in x:
            x = x.replace('this.', 'p.')
        elif 'this.' in x:
            x = x.replace('this.', 'l.')
        output_file.write(x)
        output_file.flush()

    print("Writing ...")
    input_file.close()
    output_file.close()

def fix_parser(file_path):
    print("Altering " + file_path)
    if not os.path.exists(file_path):
        print(f"Could not find file: {file_path}")
        sys.exit(1)
    parts = os.path.split(file_path)
    file_name = parts[-1]

    shutil.move(file_path, file_path + ".bak")
    input_file = open(file_path + ".bak",'r')
    output_file = open(file_path, 'w')
    for x in input_file:
        if 'this.' in x:
            x = x.replace('this.', 'p.')
        output_file.write(x)
        output_file.flush()

    print("Writing ...")
    input_file.close()
    output_file.close()

if __name__ == '__main__':
    main(sys.argv)
