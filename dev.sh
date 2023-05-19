cd ~/HyperledgerFabric/myshells/devModOn/
if (( $#==0 ))
then
    printf "give 1 param n to execute n.sh, n from 1 to 4.\nif n is 3, another param should be given to indicate the chaincode path.\n"
elif (( $#==1 ))
then
    sh $1".sh"
elif (( $#==2 && $1==3 ))
then
    sh "3.sh" $2
else
    printf "invalid parameter.\n"
fi
