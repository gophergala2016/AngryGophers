# AngryGophers - Tanks

## Description

Our version of "Tanks" is a simple 2D game "from the past", where you can drive your own tank. The idea came from the game released in 80's, where two players could fight AI opponents on one gaming console. We decided to take it a little further and created a multiplayer game for the masses. We have small battlefield, where tanks battle for supremacy. Graphic is fancy (we used Kenney assets, mostly), we created sounds by ourselves and, what is important, we had a lot of fun during development.

## Used Technologies 

No one in our team has any experience in game development, but we decided that the best way to earn it is to create our first game. Thus, without any real knowledge about upcoming challenges we tried and mostly succeded in creating this simple game. At first we wanted to simply use websockets connecting straight to the server, but after some work it looked like there would be way to much to do on the front-end side. However, Gopher Gala is all about Go, so ulimately we ditched that idea in favour of writing most of the code in Go. Both server and client are written in Go, and communication between them was switched from websockets to UDP. There is still some code written in JavaScript (drawing canvas) but the calculations are done in backend.

## Running a game

Just compile and run client script that is located in client folder. The server is running online on our side all the time so you can just join it.

## Rules

Every player gets his own tank to drive and conquer.

There are few types of ground tiles which differ in tank speed:
* dirt, grass, dirt+grass - normal type of tiles
* sand - you can't drive through it, so if you're respawned on it, there's nothing you can do but shooting around
* sand+grass - slows you down a little bit
* water - slows you down more than sand+grass
* ice - gives you extra speed
* rocks, trees - you can hide below them (yes, below the rocks too)
 
There are also some random power-ups that give you extra speed.

Players can also release the dust cloud to hide his canny intenstions.

## Key bindings

* movement          - arrows
* fire               - space
* release dust cloud - ctrl

