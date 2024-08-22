import os

os.system("sudo ./scripts/build.sh")

for subdir, dirs, files in os.walk("./examples"):
    for file in files:
        if file == "Watchfile":
            print(f"\ntesting for {subdir}\n")
            os.system(f"cd {subdir} && watch")
