import os
import re
import sys
import subprocess

cmd = "git describe"
p = subprocess.Popen(cmd.split(), stdout=subprocess.PIPE, stderr=subprocess.PIPE)
branch, err = p.communicate()

# 2.2.3-3-g2019e6e

regex = r"(\d+\.\d+\.\d+)"
if re.search(regex, branch):
  match = re.search(regex, branch)
  print match.group(0)
  sys.exit(0)

