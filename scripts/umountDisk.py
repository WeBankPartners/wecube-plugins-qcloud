#!/usr/bin/python
import sys,getopt,os,re

fstabFile="/etc/fstab"

def umountDir(mountDir):
    cmd="umount " + mountDir
    result = os.system(cmd)
    if result != 0:
        print "execute umount failed cmd=%s,result=%d" %(cmd,result)

def delMountInfoFromfstabFile(diskName,mountDir):
    with open(fstabFile,"r") as f:
        lines =f.readlines()
    with open(fstabFile,"w") as f:
        for line in lines:
            if diskName in line and mountDir in line:
	            pass
            else:
                f.write(line)
    

def main(argv):
    diskName=""
    mountDir=""

    try:
        opts,args=getopt.getopt(argv,"hd:m:",["diskName=","mountDir="])
    except getopt.GetoptError:
        print 'unmountDisk.py -d <diskName> -m <mountDir>'
        sys.exit(2)

    for opt ,arg in opts:
        if opt=='-h':
            print 'unmountDisk.py -d <diskName>  -m <mountDir>'
            sys.exit(0)
        elif opt in ("-d","--diskName"):
            diskName=arg
        elif opt in ("-m","--mountDir"):
            mountDir=arg
    
    if mountDir =="" or diskName =="":
        print "input param have some empty value"
        sys.exit(2)
    
    umountDir(mountDir)
    delMountInfoFromfstabFile(diskName,mountDir)

  
if __name__ == "__main__":
    main(sys.argv[1:])
    sys.exit(0)

