# 🧠💻Think
⚡️AI-Assisted Command Line Mastery⚡️

## Overview

`think` is a command line tool that incorporates AI to simplify and streamline your shell operations. No need to remember complex command line details, just describe what you want, and let `think` generate the shell code for you.

## Features

- **AI-powered command generation**: Just describe in your own words what you want to achieve, and Think will generate the necessary shell command(s) for you.

- **Iterative refinement**: Not satisfied with the command generated? No problem. You can provide additional feedback and guide the AI to produce a more suitable result.

- **Code editing**: Before executing, you have the chance to review and edit the proposed command to make sure it's exactly what you want.

Please note that this will send some of your data to OpenAI (your commands and the first 1000 bytes of their stdout and stderr). If you are not comfortable with this, please do not use this tool.

## Installation

Coming soon.

## Usage

Upon running, `think` opens an interactive session where you can discuss your task with the AI, review and edit the proposed code, and then execute it when ready. It has two modes, "thinking" and "running." In "thinking" mode, you can describe your task to the AI, and it will generate the necessary shell code for you. In "running" mode, you can review and edit the proposed code, and then execute it when ready. In running mode, you can hit C-c or C-d to go back to thinking mode. If the AI gets something wrong, tell it to fix it – it has access to the truncated output of the previous command, which should help it debug.

## Examples

Do you always remember the exact syntax for the `find` command? I don't. But I know what I want to do, so I can just describe it to "Think" and let it generate the command for me.

```console
think> find the 10 largest files in ~/Downloads
Find the 10 largest files in the ~/Downloads directory, sorted by size.
find ~/Downloads -type f -exec du -sh {} \; 2> /dev/null | sort -rh | head -n 10

run> find ~/Downloads -type f -exec du -sh {} \; 2> /dev/null | sort -rh | head -n 10
564M	/Users/ric/Downloads/Anaconda3-2023.03-1-MacOSX-arm64.pkg
564M	/Users/ric/Downloads/Anaconda3-2023.03-1-MacOSX-arm64 (1).pkg
224M	/Users/ric/Downloads/Skype-8.97.0.204.dmg
209M	/Users/ric/Downloads/googlechrome.dmg
191M	/Users/ric/Downloads/VSCode-darwin-universal.zip
160M	/Users/ric/Downloads/Discord.dmg
155M	/Users/ric/Downloads/Obsidian-1.1.16-universal.dmg
116M	/Users/ric/Downloads/google-cloud-cli-430.0.0-darwin-arm.tar.gz
113M	/Users/ric/Downloads/files-WhatsApp.dmg
112M	/Users/ric/Downloads/Messenger.185.0.0.5.228.dmg
```
How about some CSV wrangling?

```console
think> Please generate a simple CSV file with a few columns. At least one column should be numeric.
Creating a simple CSV file with three columns: ID, Name, Age
echo -e "ID,Name,Age\n1,Alice,30\n2,Bob,25\n3,Carol,28" > data.csv

run> echo -e "ID,Name,Age\n1,Alice,30\n2,Bob,25\n3,Carol,28" > data.csv
think> get me the first column from data.csv
Extracting the first column from data.csv
cut -d"," -f1 data.csv

run> cut -d"," -f1 data.csv
ID
1
2
3
think> please sum the third column
Calculating the sum of the third column in data.csv
awk -F"," 'NR>1 {sum += $3} END {print sum}' data.csv

run> awk -F"," 'NR>1 {sum += $3} END {print sum}' data.csv
83
```

Now let's do something more fancy:

```console
think> create the directory structure for a go project named thinktest, with a go module and git repo. Add files like README, LICENSE, .gitignore...
Creating Go project directory structure with module, Git repo, README, LICENSE, and .gitignore files.
mkdir -p thinktest/{cmd,pkg} && cd thinktest && go mod init thinktest && git init && touch README.md LICENSE .gitignore

run> mkdir -p thinktest/{cmd,pkg} && cd thinktest && go mod init thinktest && git init && touch README.md LICENSE .gitignore^C
think> the repo should live on github. My username is ryszard. The license is MIT.
Creating Go project with GitHub remote, MIT license, and adding necessary files.
mkdir -p thinktest/{cmd,pkg} && cd thinktest && go mod init thinktest && git init && echo -e "MIT License\n\nCopyright (c) $(date +%Y) $(git config user.name)\n\nPermission is hereby granted, free of charge, to any person obtaining a copy\nof this software and associated documentation files (the \"Software\"), to deal\nin the Software without restriction, including without limitation the rights\nto use, copy, modify, merge, publish, distribute, sublicense, and/or sell\ncopies of the Software, and to permit persons to whom the Software is\nfurnished to do so, subject to the following conditions:\n\nThe above copyright notice and this permission notice shall be included in all\ncopies or substantial portions of the Software.\n\nTHE SOFTWARE IS PROVIDED \"AS IS\", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR\nIMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,\nFITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE\nAUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER\nLIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,\nOUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE\nSOFTWARE." > LICENSE && touch README.md .gitignore && git add . && git commit -m "Initial commit" && git remote add origin git@github.com:ryszard/thinktest.git

run> mkdir -p thinktest/{cmd,pkg} && cd thinktest && go mod init thinktest && git init && echo -e "MIT License\n\nCopyright (c) $(date +%Y) $(git config user.name)\n\nPermission is hereby granted, free of charge, to any person obtaining a copy\nof this software and associated documentation files (the \"Software\"), to deal\nin the Software without restriction, including without limitation the rights\nto use, copy, modify, merge, publish, distribute, sublicense, and/or sell\ncopies of the Software, and to permit persons to whom the Software is\nfurnished to do so, subject to the following conditions:\n\nThe above copyright notice and this permission notice shall be included in all\ncopies or substantial portions of the Software.\n\nTHE SOFTWARE IS PROVIDED \"AS IS\", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR\nIMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,\nFITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE\nAUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER\nLIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,\nOUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE\nSOFTWARE." > LICENSE && touch README.md .gitignore && git add . && git commit -m "Initial commit" && git remote add origin git@github.com:ryszard/thinktest.git
go: creating new go.mod: module thinktest
go: to add module requirements and sums:
	go mod tidy
Initialized empty Git repository in /Users/ric/thinktest/.git/
[main (root-commit) 085e2a9] Initial commit
 4 files changed, 24 insertions(+)
 create mode 100644 .gitignore
 create mode 100644 LICENSE
 create mode 100644 README.md
 create mode 100644 go.mod
```

Note that it got the LICENSE text right, which is pretty cool. In earlier versions it was trying to download the license from the Internet, but that was obviously problematic for an AI without access to the Internet and a knowledge cutoff.

```console
think> find all the image files in ~/Downloads. copy them to ./images (create it). Anonymize the names so that I don't leak information (this is for a demo)
Finding image files in ~/Downloads, creating ./images directory, copying them with anonymized names to the directory
mkdir -p ./images && find ~/Downloads -type f \( -iname "*.jpg" -o -iname "*.jpeg" -o -iname "*.png" \) -exec sh -c 'cp -t ./images -- "$0" "$(xxd -l 8 -ps -c 8 random)"' '{}' \;

run> mkdir -p ./images && find ~/Downloads -type f \( -iname "*.jpg" -o -iname "*.jpeg" -o -iname "*.png" \) -exec sh -c 'cp -t ./images -- "$0" "$(xxd -l 8 -ps -c 8 random)"' '{}' \;
xxd: random: No such file or directory
cp: illegal option -- t
usage: cp [-R [-H | -L | -P]] [-fi | -n] [-aclpsvXx] source_file target_file
       cp [-R [-H | -L | -P]] [-fi | -n] [-aclpsvXx] source_file ... target_directory
xxd: random: No such file or directory

(plenty of those, I removed them for brevity)


cp: illegal option -- t
usage: cp [-R [-H | -L | -P]] [-fi | -n] [-aclpsvXx] source_file target_file
       cp [-R [-H | -L | -P]] [-fi | -n] [-aclpsvXx] source_file ... target_directory
think> fix it
Finding image files in ~/Downloads, creating ./images directory, copying them with anonymized names to the directory
mkdir -p ./images && find ~/Downloads -type f \( -iname "*.jpg" -o -iname "*.jpeg" -o -iname "*.png" \) -exec sh -c 'cp "$0" "./images/$(xxd -l 8 -p -u /dev/urandom).${0##*.}"' '{}' \;

run> mkdir -p ./images && find ~/Downloads -type f \( -iname "*.jpg" -o -iname "*.jpeg" -o -iname "*.png" \) -exec sh -c 'cp "$0" "./images/$(xxd -l 8 -p -u /dev/urandom).${0##*.}"' '{}' \;
think> now please give me the average dimensions of the files in images
Calculating the average dimensions of image files in the ./images directory
find ./images -type f \( -iname "*.jpg" -o -iname "*.jpeg" -o -iname "*.png" \) -exec identify -format "%wx%h\n" '{}' \; | awk -F 'x' '{total_w+=$1; total_h+=$2; count++} END {print total_w/count "x" total_h/count}'

run> find ./images -type f \( -iname "*.jpg" -o -iname "*.jpeg" -o -iname "*.png" \) -exec identify -format "%wx%h\n" '{}' \; | awk -F 'x' '{total_w+=$1; total_h+=$2; count++} END {print total_w/count "x" total_h/count}'
find: identify: No such file or directory
find: identify: No such file or directory
find: identify: No such file or directory
find: identify: No such file or directory
find: identify: No such file or directory
find: identify: No such file or directory
find: identify: No such file or directory
find: identify: No such file or directory
find: identify: No such file or directory
find: identify: No such file or directory
find: identify: No such file or directory
awk: division by zero
 source line number 1
think> fix it
Calculating the average dimensions of image files in the ./images directory using sips command
find ./images -type f \( -iname "*.jpg" -o -iname "*.jpeg" -o -iname "*.png" \) -exec sips -g pixelWidth -g pixelHeight '{}' \; | awk '/Width|Height/{sum+=$NF;count++} END {print sum/count}'

run> find ./images -type f \( -iname "*.jpg" -o -iname "*.jpeg" -o -iname "*.png" \) -exec sips -g pixelWidth -g pixelHeight '{}' \; | awk '/Width|Height/{sum+=$NF;count++} END {print sum/count}'
2029.41
think> now please give me the average dimensions (height and width) of the files in images
Calculating the average dimensions (width and height) of image files in the ./images directory using sips command
find ./images -type f \( -iname "*.jpg" -o -iname "*.jpeg" -o -iname "*.png" \) -exec sips -g pixelWidth -g pixelHeight '{}' \; | awk '/Width/{w+=$NF;wc++} /Height/{h+=$NF;hc++} END {print "Width: " w/wc ", Height: " h/hc}'

run> find ./images -type f \( -iname "*.jpg" -o -iname "*.jpeg" -o -iname "*.png" \) -exec sips -g pixelWidth -g pixelHeight '{}' \; | awk '/Width/{w+=$NF;wc++} /Height/{h+=$NF;hc++} END {print "Width: " w/wc ", Height: " h/hc}'
Width: 2644.09, Height: 1414.73
think> now please resize all of them so that they are 1/4 of the size of the average dimensions.
Resizing all image files in ./images to 1/4 of the average dimensions (width and height) calculated previously
find ./images -type f \( -iname "*.jpg" -o -iname "*.jpeg" -o -iname "*.png" \) -exec sips -Z $((2644/4))x$((1415/4)) '{}' \;

run> find ./images -type f \( -iname "*.jpg" -o -iname "*.jpeg" -o -iname "*.png" \) -exec sips -Z $((2644/4))x$((1415/4)) '{}' \;
images/690EC1F13FB62A19.jpeg
  ./images/690EC1F13FB62A19.jpeg
./images/E0E26311A30C6376.png
  ./images/E0E26311A30C6376.png
./images/BCB68ED864C9F4F2.jpg
  ./images/BCB68ED864C9F4F2.jpg
./images/D87454DB4AA998EC.jpg
  ./images/D87454DB4AA998EC.jpg
./images/CCDE085970025EA2.png
  ./images/CCDE085970025EA2.png
./images/26DD89AF3D7B67EE.png
  ./images/26DD89AF3D7B67EE.png
./images/8A2406199B92A538.png
  ./images/8A2406199B92A538.png
./images/D705AA10B2A43BB9.jpg
  ./images/D705AA10B2A43BB9.jpg
./images/B18EB23B54B98725.png
  ./images/B18EB23B54B98725.png
./images/26D909E004CAC1C1.png
  ./images/26D909E004CAC1C1.png
./images/B7590602F4CEFE9B.jpg
  ./images/B7590602F4CEFE9B.jpg
think> show me their dimensions
Showing the dimensions of the resized image files in the ./images directory using sips command
find ./images -type f \( -iname "*.jpg" -o -iname "*.jpeg" -o -iname "*.png" \) -exec sips -g pixelWidth -g pixelHeight '{}' \;

run> find ./images -type f \( -iname "*.jpg" -o -iname "*.jpeg" -o -iname "*.png" \) -exec sips -g pixelWidth -g pixelHeight '{}' \;
./images/690EC1F13FB62A19.jpeg
  pixelWidth: 484
  pixelHeight: 661
./images/E0E26311A30C6376.png
  pixelWidth: 661
  pixelHeight: 452
./images/BCB68ED864C9F4F2.jpg
  pixelWidth: 440
  pixelHeight: 661
./images/D87454DB4AA998EC.jpg
  pixelWidth: 661
  pixelHeight: 459
./images/CCDE085970025EA2.png
  pixelWidth: 661
  pixelHeight: 661
./images/26DD89AF3D7B67EE.png
  pixelWidth: 661
  pixelHeight: 184
./images/8A2406199B92A538.png
  pixelWidth: 661
  pixelHeight: 661
./images/D705AA10B2A43BB9.jpg
  pixelWidth: 661
  pixelHeight: 372
./images/B18EB23B54B98725.png
  pixelWidth: 661
  pixelHeight: 661
./images/26D909E004CAC1C1.png
  pixelWidth: 661
  pixelHeight: 184
./images/B7590602F4CEFE9B.jpg
  pixelWidth: 495
  pixelHeight: 661
```

## Contributing

Just send me a pull request.

## License

