# goSorl
Golangs SORL : Solution ORchestration Language 

#
# SORL Scripting example
#
wait ?||$||#<br>
cat /etc/passwd<br>
var a=a11<br>
var b=b22<br>
tag test1 { <br>
   echo "{b}{b}==>{a}"<br>
   date<br>
   hostname<br>
   echo "Done tag test1"<br>
}<br>
echo "{a}{a}=>{b}" > /tmp/dellme.dell<br>
cat /tmp/dellme.dell<br>
echo "***DONE***"<br>
#<br>
<br>
tag abcd  {<br>  
   cd /tmp<br>
   pwd<br>
   echo "Done tag abcd."<br>
}<br>
<br>
if true {<br>
   echo ""<br>
   echo "*** IF works fine***"<br>
   echo ""<br>
}<br>
echo "SORL DONE"<br>
exit<br>
<br>
