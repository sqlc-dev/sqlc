import dataclasses
from typing import Optional

# Comment
# Longer Comment

@dataclasses.dataclass()
class InventoryItem:
    """Class for keeping track of an item in inventory."""
    name: Optional[str] = None
    unit_price: Optional[float] = None
    quantity_on_hand: Optional[int] = None

    def total_cost(self) -> float:
        return self.unit_price * self.quantity_on_hand

if __name__ == "__main__":
    print(InventoryItem())

