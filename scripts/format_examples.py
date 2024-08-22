import os

tokens = [
    "CLONE",
    "SET",
    "EXTRACT",
    "RUN"
]

def format_file(subdir, file):
    print(f"\nformatting {subdir}/{file}\n")
    reader = open(f"{subdir}/{file}", "r")
    lines = reader.readlines()
    reader.close()

    newContent = ""
    for l in lines:
        anyFound = False
        for tok in tokens:
            if l.startswith(tok) and l[len(tok)] == "(":
                newContent += l[:len(tok)] + " " + l[len(tok):]
                anyFound = True
        if not anyFound:
            newContent += l


    writer = open(f"{subdir}/{file}", "w")
    writer.write(newContent)
    writer.close()

for subdir, dirs, files in os.walk("./examples"):
    for file in files:
        if file == "Watchfile":
            format_file(subdir, file)
