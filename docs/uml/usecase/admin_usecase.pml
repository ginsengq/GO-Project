@startuml

actor Admin

Admin --> (Manage Users)
Admin --> (Manage Cars)
Admin --> (Manage Orders)
Admin --> (Manage Transactions)

(Manage Users) --> (Add User)
(Manage Users) --> (Edit User)
(Manage Users) --> (Delete User)
(Manage Cars) --> (Add Car)
(Manage Cars) --> (Edit Car)
(Manage Cars) --> (Delete Car)
(Manage Cars) --> (Change Car Status)
(Manage Orders) --> (View Order)
(Manage Orders) --> (Confirm Order)
(Manage Orders) --> (Cancel Order)
(Manage Transactions) --> (View Transactions)
(Manage Transactions) --> (Process Payment)

@enduml
