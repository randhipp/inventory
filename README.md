


## Task 1 : Online Store

#### ### Q1: Describe what you think happened that caused those bad reviews during our 12.12 event and why it happened.

##### Answer :

1. You got a race condition on the checkout and payment flow.
2. Maybe, No stock check before user do a payment.
3. Qty should not have a negative value, although some business allow this behaviour. 

( On my previous company, we can have negative stock so customer still can do a payment and company not loosing some sales. Company can order it and said to customer will delivery it 2-3 days again )

#### ### Q2. Based on your analysis, propose a solution that will prevent the incidents from occurring again.

##### Answer :

It is a quite challenging problem, and should be a business problem not a software problem.

https://softwareengineering.stackexchange.com/questions/412515/should-an-e-commerce-application-reserve-products-before-attempting-payment

On Software Engineer POV, we can do things like:

1. Add a reserved system on an item
2. Auto change transaction for the losing customer to pre-order ( of course with terms updated/stated in the company t&c )
3. Check all code, for race condition with go race detector and fix that 
```
go run --race race.go
```

#### Option 1 : Add a reserve system on an item
```
  Condition
  - Item A with stock qty is => 1 pc
  - User 1, 2, 3 add to chart in the same time 0.1-1ms ahead of each
```
  The business flow and/or terms of use should mention :
   - The stock will not be guaranteed to you until you pay
 
  User 1 POV
  - When user 1 add an item to cart, deduct the quantity of that item, and had a new table for reserved_item and set expired time for payment
  - If user 1 checkout, do the payment, the item will belongs to user 1 => item A will be 0
  - If user 1 not continue the checkout and not doing any payment, after the expired time is passing, the item will back available.

  User 2 POV
  - user 2 still got the item to cart
  - will fail at payment flow ( before payment can continue, we check again the qty )

  User 3 POV
  - user 3 will not success add item to cart

  Pro : 
  - Stock quantity will be match and no race condition
  
  Cons :
  - Bad User Experience on user 2, they will asume they got the item, but fail when doing some payment.
  - If the expired time is too long and user 1 abandoned the cart, we may lose sales from user 2.
  - Business will lose sales, because item can't be back ordered.


#### Option 2 : Auto change transaction for the losing customer to pre-order ( of course with terms updated/stated in the company t&c )



Do a locking mechanism
- client add item to cart