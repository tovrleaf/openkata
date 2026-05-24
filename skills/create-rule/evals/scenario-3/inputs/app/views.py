from django.db import connection
from django.shortcuts import get_object_or_404
from .models import Order, Product, Customer

def get_order_summary(request, customer_id):
    # Raw SQL for performance - this bypasses the ORM
    with connection.cursor() as cursor:
        cursor.execute(
            "SELECT o.id, o.total, o.created_at FROM orders_order o "
            "WHERE o.customer_id = %s ORDER BY o.created_at DESC",
            [customer_id]
        )
        rows = cursor.fetchall()
    return rows

def get_active_products(request):
    # N+1 query: iterating and calling .category for each product
    products = Product.objects.filter(active=True)
    result = []
    for p in products:
        result.append({
            "name": p.name,
            "category": p.category.name,  # separate query per product
            "price": p.price,
        })
    return result

def get_customer_orders(customer_id):
    customer = get_object_or_404(Customer, pk=customer_id)
    # Uses ORM but no select_related
    orders = Order.objects.filter(customer=customer)
    return list(orders.values("id", "total", "status"))
