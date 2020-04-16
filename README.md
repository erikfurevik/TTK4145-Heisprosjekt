This is our elevator system made for the elevator project in the course TTK4145 Real-Time Programming at NTNU. 

The network module used in this project was handed out by the course and can be found [here](https://github.com/TTK4145/Network-go). This network module is made for the programming language Google GO.

We used a peer-to-peer solution where every elevator has all information about all elevators and calculate the costfunction for every elevator. This way all elevators will agree which elevator should take what orders. Since all elevators know all information about all orders and the status of the other elevators it is easy to handle errors like one elevator going offline or a motor shutting down etc. 