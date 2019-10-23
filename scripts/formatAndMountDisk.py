#!/usr/bin/python

import sys,getopt,os,re

def execCmd(cmd):
    r=os.popen(cmd)
    text=r.read()
    r.close()
    return text

def ensureDir(mountDir):
    if not os.path.exists(mountDir):
        os.makedirs(mountDir)

def isFormatedDisk(diskName):
    result=execCmd("lsblk -p -P -o  FSTYPE "+diskName)
    rows=result.split("\n")
    for row in rows:
        fsType=row[len("FSTYPE=")+1:-1]
        if len(fsType) > 0:
            return True
    return False       

def mountDisk(diskName,mountDir,fileSystemType):
    mountOptions={
        "ext3":"    noatime,acl,user_xattr 1 1",
        "ext4":"    noatime,acl,user_xattr 1 1",
        "xfs":"    defaults     0    2",
    }

    cmd="mount " + diskName+" " + mountDir
    result=os.system(cmd)
    if result != 0:
        print "mount cmd exec failed:cmd=%s result=%d" % (cmd,result)
        sys.exit(1)
    
    fstab=diskName+"        "+mountDir+"        "+fileSystemType+"        "+mountOptions[fileSystemType]
    with open("/etc/fstab","a") as f:
        f.write(fstab) 


def formatDisk(diskName,fileSystemType):
    formatDict={
        "ext3":"mkfs.ext3 -F ",
        "ext4":"mkfs.ext4 -F ",
        "xfs": "mkfs.xfs -n ftype=1  -f "
    }

    cmd=formatDict[fileSystemType] + diskName
    result = os.system(cmd)
    if result != 0:
        print "execute format disk failed cmd=%s,result=%d" %(cmd,resullt)
        sys.exit(1)

def main(argv):
    diskName=""
    mountDir=""
    fileSystemType=""

    try:
        opts,args=getopt.getopt(argv,"hd:f:m:",["diskName=","fileSystemType=","mountDir="])
    except getopt.GetoptError:
        print 'formatAndMountDisk.py -d <diskName> -f <formatFileSystemType> -m <mountDir>'
        sys.exit(2)

    for opt ,arg in opts:
        if opt=='-h':
            print 'formatAndMountDisk.py -d <diskName> -f <formatFileSystemType> -m <mountDir>'
            sys.exit(0)
        elif opt in ("-d","--diskName>"):
            diskName=arg
        elif opt in ("-f","--formatFileSystemType"):
            fileSystemType=arg
        elif opt in ("-m","--mountDir"):
            mountDir=arg
    
    if mountDir =="" or diskName =="" or fileSystemType == "":
        print "input param have some empty value"
        sys.exit(2)
    
    if fileSystemType not in ("ext3","ext4","xfs"):
        print "invalid fileSystemType(%s)" % fileSystemType
        sys.exit(2)

    isFormated = isFormatedDisk(diskName)
    if isFormated:
        print "disk(%s) has been formated" % diskName
        sys.exit(1)

    ensureDir(mountDir)
    formatDisk(diskName,fileSystemType)
    mountDisk(diskName,mountDir,fileSystemType)
  
if __name__ == "__main__":
    main(sys.argv[1:])
    sys.exit(0)



