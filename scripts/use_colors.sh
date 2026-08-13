RED="\x1b[0;31m"
YELLOW="\x1b[0;33m"
CYAN="\x1b[0;36m"
GREEN="\x1b[0;32m"
BLUE="\x1b[0;34m"
PURPLE="\x1b[0;35m"
RESET="\x1b[0;0m"

function red() {
    input=$1
    if [[ -z $input ]] && [[ ! -t 0 ]]; then
        while read -r line ; do
            input+=$line
        done
    fi

    echo -e "${RED}${input}${RESET}"
}
export -f red

function yellow() {
    input=$1
    if [[ -z $input ]] && [[ ! -t 0 ]]; then
        while read -r line ; do
            input+=$line
        done
    fi

    echo -e "${YELLOW}${input}${RESET}"
}
export -f yellow

function cyan(){
    input=$1
    if [[ -z $input ]] && [[ ! -t 0 ]]; then
        while read -r line ; do
            input+=$line
        done
    fi

    echo -e "${CYAN}${input}${RESET}"
}
export -f cyan

function green(){
    input=$1
    if [[ -z $input ]] && [[ ! -t 0 ]]; then
        while read -r line ; do
            input+=$line
        done
    fi

    echo -e "${GREEN}${input}${RESET}"
}
export -f green

function blue(){
    input=$1
    if [[ -z $input ]] && [[ ! -t 0 ]]; then
        while read -r line ; do
            input+=$line
        done
    fi

    echo -e "${BLUE}${input}${RESET}"
}
export -f blue

function purple(){
    input=$1
    if [[ -z $input ]] && [[ ! -t 0 ]]; then
        while read -r line ; do
            input+=$line
        done
    fi

    echo -e "${PURPLE}${input}${RESET}"
}
export -f purple

function printcolor(){
    input=$1
    if [[ -z $input ]] && [[ ! -t 0 ]]; then
        while IFS= read -r line ; do
            printcolor "$line"
        done

        return 0
    fi

    echo -e "$input"

    return 0
}
export -f printcolor

function fileToRedMsg() {
    input="$1"
    if [[ -z $input ]] && [[ ! -t 0 ]]; then
        while IFS='\n' read -r line ; do
            fileToRedMsg "$line"
        done

        return 0
    fi

    output="$input"
    output=$(echo "$output" | sed -r "s/([\\_A-Za-z0-9\\/\\.\\-]+:[0-9]+)/${RED}\\1${RESET}/")

    printcolor "$output"
}
export -f fileToRedMsg

if [ ! -t 0 ]; then
    while read -r line ; do
        echo $line | printcolor
    done
fi
