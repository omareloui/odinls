# What I need from this app

## Near future

- calculate the product's price from the time it takes and the materials cost
- print the order's invoice
- fix the bugs with authorization
- make the order updatable
- make a page to view the order status (and a page for the craftsmen to update
  the status)
- know how to navigate the app as a user, not a craftsman
- list the products for the end user
- move the current routes to `/dashboard` parent directory

### Bugs

- on creating the order get the new created order (the variant id doesn't show-up)
- updating the order doesn't work

## Far future

- add inventory
- connect the inventory with the product and calculate the total products i can
  do with my current inventory
- design
- self healing links for the products

### Materials

Implement this to make the material cost dynamically calculable.

```go
type Materials struct {
  ID: primitive.ObjectID;
  Kind: string; // enum: "leather" | "hardware" | "fabric" | "dye"
  Supplier: primitive.ObjectID;
  PricePerUnit: float64;
  Unit: string; // enum: "ft²" | "unit" | "piece" | ...
  InventoryCount: float64;
}

type Supplier struct {
  ID: primitive.ObjectID;
  Name: string;
  Location: string;
  Notes: string;
}

type Variant struct {
  // ...
  RequiredMaterilas struct {
    material: primitive.ObjectID;
    Units: float64;
  }
}

func (v *Variant) CalculateMaterialCost() float64
```
