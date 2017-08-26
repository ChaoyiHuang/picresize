# picresize to demonstrate edge cloud computing as a service, based on Kubernetes

picserver:
  a web server to allow picture uploading, and start "picresize" container at specified edge node through the kubernetes master in the cloud, the picresize container will resize the uploaded pictures with two small size ones: the width is 128 and 256, these two small ones will be pushed back to the picserver in the cloud.

  all edge nodes which could be behind firewall and/or NAT can register to the kubernetes master in the cloud.

picresize:
  a container running in the edge node, pull the picture from the picserver, and resize these picture with two small ones, then upload the small ones to the picserver. picresize container will exit after the resizing job is finished. picresize must be compiled in static, so that the container could be based on alpine base image.

***ALL SOURCE CODE IS WRITEN IN HARD CODE TO MAKE THE DEMO WORK SOON***
