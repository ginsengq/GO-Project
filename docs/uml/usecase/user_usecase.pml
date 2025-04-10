@startuml

actor User

User --> (Browse Catalog)
User --> (Reserve Car)
User --> (Order Car)
User --> (Schedule Test Drive)
User --> (Make Deposit)
User --> (Complete Order)

(Browse Catalog) --> (Browse Car Details)
(Reserve Car) --> (Select Car)
(Reserve Car) --> (Confirm Reservation)
(Order Car) --> (Select Car)
(Order Car) --> (Make Payment)
(Order Car) --> (Confirm Order)
(Schedule Test Drive) --> (Choose Car)
(Schedule Test Drive) --> (Select Date)
(Make Deposit) --> (Choose Amount)
(Make Deposit) --> (Confirm Payment)
(Complete Order) --> (Confirm Payment)

@enduml
