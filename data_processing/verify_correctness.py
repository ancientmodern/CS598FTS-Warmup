
import os

def isT1BeforeT2(t1: list, t2: list):
    if t1[0] < t2[0]: return True
    elif t1[0] > t2[0]: return False
    else:
        if t1[1] <= t2[1]: return True
        else: return False

directory = "logs_for_correctness/logs_with_32_clients"
 
for filename in os.listdir(directory):
    data = {}
    filepath = os.path.join(directory, filename)
    if not os.path.isfile(filepath):
        print("not a file {}".format(filepath))
        continue
    f = open(filepath, 'r')
    for line in f.readlines():
        line = line.split()
        if line[7] not in data:
            data[line[7]] = [int(line[9]), int(line[11])]
        else:
            if isT1BeforeT2(data[line[7]], [int(line[9]), int(line[11])]): data[line[7]] = [int(line[9]), int(line[11])]
            else:   
                print("NOT linear!")
                print(filename,  "line: ", line)
                f.close()
                exit(1)
    f.close()

print("Linearizablity Verified!")


