DIR=$(date)

mkdir "$DIR"
mv *.{caffemodel,solverstate} "$DIR"
cp *{txt,py} "$DIR"

