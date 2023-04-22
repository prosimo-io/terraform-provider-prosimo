## This script is used for auto versioning of terraform binary.
#!/bin/bash

file=Makefile
input_flag=$1
current_major_version=`sed -n 's/^ *MAJOR_VERSION*= *//p' "$file"`
current_minor_version=`sed -n 's/^ *MINOR_VERSION*= *//p' "$file"`
current_revision=`sed -n 's/^ *REVISION*= *//p' "$file"`
# current_build_number = `sed -n 's/^ *BUILDNUMBER*= *//p' "$file"`

if [ "$input_flag" == "major-version" ]; then
    current_major_version=$(( $current_major_version + 1 ))
elif [ "$input_flag" == "minor-version" ]; then
    current_minor_version=$(( $current_minor_version + 1 ))
elif [ "$input_flag" == "revision" ]; then
    current_revision=$(( $current_revision + 1 ))
# elif [ "$input_flag" == "build_number" ]; then
#     current_build_number=$(( $current_build_number + 1 ))
else
    continue
fi



sed -i '' 's/MINOR_VERSION=.*/MINOR_VERSION='$current_minor_version'/' "$file"
sed -i '' 's/MAJOR_VERSION=.*/MAJOR_VERSION='$current_major_version'/' "$file"
sed -i '' 's/REVISION=.*/REVISION='$current_revision'/' "$file"