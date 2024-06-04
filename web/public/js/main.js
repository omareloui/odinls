function parseMoney(str) {
  const numStr = str.replace(/\D/g, "");
  if (!numStr) {
    return NaN;
  }
  return parseInt(numStr);
}

function calculateItemPrice(products, item) {
  const customPrice = item.custom_price?.value;
  if (customPrice) {
    const num = parseMoney(customPrice);
    if (!Number.isNaN(num)) {
      return num;
    }
  }
  const prodId = item.product_id?.value;
  const varId = item.variant_id?.value;
  if (!varId || !prodId) return 0;
  const prod = products.find((x) => x.id === prodId);
  if (!prod) return 0;
  const variant = prod.variants.find((x) => x.id === varId);
  return variant?.price || 0;
}

function calculateItemTotal(products, item) {
  const itemPrice = calculateItemPrice(products, item);
  const quantity = parseInt(item.quantity?.value || "1");
  return itemPrice * quantity;
}

function calculateSubtotal(products, items, priceAddons) {
  let sum = 0;

  items.forEach((item) => (sum += calculateItemTotal(products, item)));

  let discountsPercentage = 0;
  let taxesPercentage = 0;
  let feesPercentage = 0;
  let shippingPercentage = 0;

  priceAddons.forEach((addon) => {
    const amount = parseMoney(addon.amount?.value || "");

    if (!addon.kind?.value || !amount || Number.isNaN(amount)) {
      return;
    }

    if (!addon.is_percentage?.value) {
      if (addon.kind?.value === "discount") {
        return (sum -= amount);
      }
      return (sum += amount);
    }
    switch (addon.kind?.value) {
      case "fees": {
        feesPercentage += amount / 100;
        break;
      }
      case "taxes": {
        taxesPercentage += amount / 100;
        break;
      }
      case "shipping": {
        shippingPercentage += amount / 100;
        break;
      }
      case "discount": {
        discountsPercentage += amount / 100;
        break;
      }
    }
    sum += sum * feesPercentage;
    sum += sum * shippingPercentage;
    sum -= sum * discountsPercentage;
    sum += sum * taxesPercentage;
  });

  return sum;
}
