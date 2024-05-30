import { ref } from 'vue';

export const cart = ref({
  items: [],
  total: "0.00",
  empty() {
    this.items = [];
    this.total = "0.00";
  },
  getCount(id) {
    let found = this.items.filter((item) => item.id === id);
    if (found && found.length === 1) {
      return found[0].count;
    }

    return 0;
  },
  addItem(id, name, price) {
    let found = false;
    this.items.every((item, index) => {
      if (item.id === id) {
        found = true;
        let temp = this.items[index];
        temp.count += 1;
        temp.price = (parseFloat(temp.price) + price).toFixed(2);
        this.items[index] = temp;
        return false;
      }

      return true;
    });

    if (found === false) {
      this.items.push({
        id: id,
        name: name,
        count: 1,
        price: price.toFixed(2),
      });
    }

    let temp = parseFloat(this.total);
    temp += price;
    this.total = temp.toFixed(2);
  },
  removeItem(id, price) {
    let found = false;
    let deleteIndex = null;
    this.items.every((item, index) => {
      if (item.id === id) {
        found = true;
        let temp = this.items[index];

        if (temp.count === 1) {
          deleteIndex = index;
          return false;
        }

        temp.count -= 1;
        temp.price = (parseFloat(temp.price) - price).toFixed(2);
        this.items[index] = temp;
        return false;
      }

      return true;
    });

    if (deleteIndex != null) {
      this.items = this.items.filter((item) => item.id != id);
    }

    let temp = parseFloat(this.total);
    temp -= price;
    this.total = temp.toFixed(2);
  },
});
