import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { Table, Button, message, Empty } from 'antd';

function Cart() {
    const [cartItems, setCartItems] = useState([]);
    const [loading, setLoading] = useState(true); // Для отображения загрузки

    useEffect(() => {
        fetchCartItems();
    }, []);

    const fetchCartItems = () => {
        setLoading(true); // Включаем индикатор загрузки
        axios.get('/api/cart/items')
            .then(res => {
                setCartItems(res.data || []);
                setLoading(false); // Отключаем индикатор загрузки
            })
            .catch(err => {
                console.error(err);
                setLoading(false); // Отключаем индикатор загрузки
            });
    };

    const handleRemove = (itemId) => {
        axios.delete(`/api/cart/items/${itemId}`)
            .then(() => {
                message.success('Item removed from cart');
                fetchCartItems();
            })
            .catch(err => console.error(err));
    };

    const handleDecreaseQuantity = (itemId) => {
        axios.patch(`/api/cart/items/${itemId}/decrease`)
            .then(() => {
                message.success('Item quantity decreased');
                fetchCartItems(); // Обновляем корзину
            })
            .catch(err => {
                console.error(err);
                message.error('Failed to decrease item quantity');
            });
    };

    const handleIncreaseQuantity = (itemId, currentQuantity, stock) => {
        if (currentQuantity >= stock) {
            message.error('Cannot add more items than available in stock');
            return;
        }
    
        axios.patch(`/api/cart/items/${itemId}/increase`)
            .then(() => {
                message.success('Item quantity increased');
                fetchCartItems();
            })
            .catch(err => {
                console.error(err);
                message.error('Cannot add more items than available in stock');
            });
    };
    

    const handleCheckout = () => {
        axios.post('/api/orders/checkout')
            .then((res) => {
                message.success(`Order created successfully. Order ID: ${res.data.order_id}`);
                setCartItems([]); // Очищаем корзину на фронте
            })
            .catch((err) => {
                console.error('Error during checkout:', err);
                message.error('Failed to complete checkout');
            });
    };

    const columns = [
        { title: 'Product', dataIndex: 'product_name', key: 'product_name' },
        { title: 'Quantity', dataIndex: 'quantity', key: 'quantity' },
        { title: 'Price', dataIndex: 'price', key: 'price' },
        {
            title: 'Action',
            key: 'action',
            render: (_, record) => (
                <>
                    <Button onClick={() => handleDecreaseQuantity(record.cart_item_id)} disabled={record.quantity <= 1}>
                        -
                    </Button>
                    <Button onClick={() => handleIncreaseQuantity(record.cart_item_id)} style={{ marginLeft: 8 }} disabled={record.quantity >= record.stock}>
                        +
                    </Button>
                    <Button onClick={() => handleRemove(record.cart_item_id)} style={{ marginLeft: 8 }}>
                        Remove
                    </Button>
                </>
            ),
        },
    ];

    return (
        <div>
            <h2>Your Cart</h2>
            {cartItems.length === 0 && !loading ? ( // Если корзина пуста и данные загружены
                <Empty description="Your cart is empty" />
            ) : (
                <Table
                    dataSource={cartItems}
                    columns={columns}
                    rowKey="cart_item_id"
                    loading={loading} // Показать индикатор загрузки
                />
            )}
            <Button
                type="primary"
                onClick={handleCheckout}
                disabled={cartItems.length === 0} // Блокируем кнопку, если корзина пуста
                style={{ marginTop: 16 }}
            >
                Proceed to Checkout
            </Button>
        </div>
    );
}

export default Cart;
