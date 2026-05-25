from django.db import models

class Category(models.Model):
    name = models.CharField(max_length=100)

    class Meta:
        db_table = "catalog_category"

class Product(models.Model):
    name = models.CharField(max_length=200)
    price = models.DecimalField(max_digits=10, decimal_places=2)
    active = models.BooleanField(default=True)
    category = models.ForeignKey(Category, on_delete=models.CASCADE)

    class Meta:
        db_table = "catalog_product"

class Customer(models.Model):
    email = models.EmailField(unique=True)
    name = models.CharField(max_length=200)

    class Meta:
        db_table = "customers_customer"

class Order(models.Model):
    STATUS_CHOICES = [("pending", "Pending"), ("paid", "Paid"), ("shipped", "Shipped")]
    customer = models.ForeignKey(Customer, on_delete=models.CASCADE)
    total = models.DecimalField(max_digits=12, decimal_places=2)
    status = models.CharField(max_length=20, choices=STATUS_CHOICES)
    created_at = models.DateTimeField(auto_now_add=True)

    class Meta:
        db_table = "orders_order"
