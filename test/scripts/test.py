import sys 

def print_parameters(params):
    print("COUNT PARAMS:", len(params))
    print("PARAMS:")
    for param in params:
        print(param)

print("PYTHON")
print_parameters(sys.argv[1:])
