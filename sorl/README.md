# goSorl
Golangs SORL : Solution ORchestration Language 

#
# SORL Scripting example
#
wait ?||$||#
cat /etc/passwd
var a=a11
var b=b22
tag test1 { 
   echo "{b}{b}==>{a}"
   date
   hostname
   echo "Done tag test1"
}
echo "{a}{a}=>{b}" > /tmp/dellme.dell
cat /tmp/dellme.dell
echo "***DONE***"
#

tag abcd  {  
   cd /tmp
   pwd
   echo "Done tag abcd."
}

if true {
   echo ""
   echo "*** IF works fine****"
   echo ""
}
echo "SORL DONE"
exit
