
https://sessionize.com/app/speaker/session/invite/577415/ciouvya


https://rentola.nl/en/listings/julianaplein-7-frankendael-1097dn-amsterdam-bc6773

https://www.kamer.nl/en/rent/apartment/amsterdam/512533/gustav-mahlerlaan/

https://www.kamer.nl/en/rent/apartment/amsterdam/463279/albert-cuypstraat/

https://www.kamer.nl/en/rent/apartment/amsterdam/467594/albert-cuypstraat/

# picresize to demonstrate edge cloud computing as a service, based on Kubernetes

***ALL SOURCE CODE IS HARD CODE TO MAKE THE DEMO/IDEA WORK ASAP***

picserver:
  a web server to allow picture uploading, and start "picresize" container at specified edge node
  through the kubernetes master in the   cloud, the picresize container will resize the uploaded
  pictures with two small size ones: the width is 128 and 256, these two small ones will be pushed
  back to the picserver in the cloud.
  
  all edge nodes which could be behind firewall and/or NAT can register to the kubernetes master
  in the cloud.

picresize:
  a container running in the edge node, pull the picture from the picserver, and resize the 
  picture with two small ones, then upload the small ones to the picserver. picresize container
  will exit after the resizing job is finished. picresize must be compiled linking staticly, so that
  the container could work using alpine as base image.

    CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags '-s' -o picresize picresize.go
  
  During the demo, picresize container image will be pulled from Docker hub if it's not in the
  edge node for the first time. After that, local image will be used instead. You can also pull
  the image manually:

     docker pull joehuang/picresize

How to register an edge node to the server
  *********
