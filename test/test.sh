for frame in 10 100 500;
do
    for thread in "0 2 4 6 8 ";
    do
        for i in `seq 1 5`;
        do
            (time ../main/main -frame $frame -p=true -threads $thread -resolution large) 2> csv_$thread\_$frame\_$i.txt        
         done
    done
done


for frame in 10 100 500;
do
    for i in `seq 1 5`;
    do
        (time ../main/main -frame $frame -resolution large) 2> csv_$thread\_seq\_$i.txt        
    done
done