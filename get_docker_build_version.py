import os
import re
import sys
import subprocess

cmd = "git branch"
p = subprocess.Popen(cmd.split(), stdout=subprocess.PIPE, stderr=subprocess.PIPE)
branch, err = p.communicate()

regex = r"release/(.*)\n"
if re.search(regex, branch):
  match = re.search(regex, branch)
  print match.group(0)
  sys.exit(0)

cmd = "git rev-parse --short HEAD"
p = subprocess.Popen(cmd.split(), stdout=subprocess.PIPE, stderr=subprocess.PIPE)
shortSha, err = p.communicate()
print shortSha.replace('\n', '')
sys.exit(0)

