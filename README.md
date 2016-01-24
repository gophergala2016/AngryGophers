# AngryGophers - Tanks

## Description

"Tanks" is a simple game "from the old days". We have small battlefield, where tanks could battle for supremacy.

## Used Technologies 

No one in our team has any experience in game development, but we decided that the best way to learn is to do. Thus without any real knowledge about comming challenges we tried, and mostly succeded in creating this simple game. At first we wanted to simply use websockets connecting straight to the server, but after some work it looked like there would be way to much to do on the front-end side. Gopher gala is all about Golang however, so ulimately we ditched that idea in favour of writing most of the code in Go. Both server and client are written in Go, and communication between them was switched from websockets to UDP. There is still some code written in Javascript but most of the calculations is done in backend.

## Running a game

Just compile and run client script that iss located in client folder.

## key bindings
movement          - arrows

fire              - space

create dust cloud - ctrl

