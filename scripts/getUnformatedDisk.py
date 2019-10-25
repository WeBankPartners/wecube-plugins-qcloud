#!/usr/bin/python

import os
import sys
import json

def execCmd(cmd):
    r=os.popen(cmd)
    text=r.read()
    r.close()
    return text

def getDiskTypeAndNameFromRow(row):
    infos=row.split(" ")
    if len(infos) != 2:
        return "","" 

    return infos[0][len("TYPE=")+1:-1],infos[1][len("NAME=")+1:-1]

def getHostDisks():
    disks=[]
    result=execCmd("lsblk -d -P -p -o TYPE,NAME")
    rows=result.split("\n")
    for row in rows:
        diskType,name=getDiskTypeAndNameFromRow(row)
        if diskType == "disk" and len(name) > 0 :
            disks.append(name)    

    return disks

def isUnformatedDisk(diskName):
    result=execCmd("lsblk -p -P -o  FSTYPE "+diskName)
    rows=result.split("\n")
    for row in rows:
        fsType=row[len("FSTYPE=")+1:-1]
        if len(fsType) > 0:
            return False
    return True        
   
if  __name__=='__main__':
    unFormatedDiskList=[]
    diskList = getHostDisks()
    for disk in diskList:
        isUnformated=isUnformatedDisk(disk)
        if isUnformated:
            unFormatedDiskList.append(disk)
   
    dict={"unformatedDisks":unFormatedDiskList}
    print json.dumps(dict)
    sys.exit(0)
