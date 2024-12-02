import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import axios from 'axios';
import { Table, Button, message, Modal, Select } from 'antd';

const { Option } = Select;

function OrderDetails() {
    const { id } = useParams(); // Получаем ID заказа из URL
    const navigate = useNavigate();
    const [orderItems, setOrderItems] = useState([]); // Инициализация пустым массивом
    const [loading, setLoading] = useState(true);
    const [isPaid, setIsPaid] = useState(false);
    const [paymentModalVisible, setPaymentModalVisible] = useState(false);
    const [selectedPaymentMethod, setSelectedPaymentMethod] = useState('');

    useEffect(() => {
        setLoading(true);
        axios.get(`/api/orders/${id}/items`) // Запрос на получение элементов заказа
            .then(res => {
                setOrderItems(res.data.items || []); // Убедитесь, что items — это массив
                setIsPaid(res.data.status === 'Paid');
                setLoading(false);
            })
            .catch(err => {
                console.error(err);
                setOrderItems([]); // На случай ошибки сбрасываем в пустой массив
                setLoading(false);
            });
    }, [id]);

    const openPaymentModal = () => setPaymentModalVisible(true);
    const closePaymentModal = () => setPaymentModalVisible(false);

    const handlePayment = () => {
        if (!selectedPaymentMethod) {
            message.error('Please select a payment method');
            return;
        }

        const totalAmount = orderItems.reduce(
            (sum, item) => sum + item.price * item.quantity,
            0
        );
        
        const order_id = parseInt(id, 10); 
        axios.post(`/api/payments`, {
            order_id: order_id,
            amount: totalAmount,
            payment_method: selectedPaymentMethod,
        })
            .then(() => {
                message.success('Payment recorded successfully');
                setIsPaid(true);
                closePaymentModal();
            })
            .catch(err => {
                console.error('Error recording payment:', err);
                message.error('Failed to record payment');
            });
    };

    const columns = [
        { title: 'Product Name', dataIndex: 'product_name', key: 'product_name' },
        { title: 'Product ID', dataIndex: 'product_id', key: 'product_id' },
        { title: 'Quantity', dataIndex: 'quantity', key: 'quantity' },
        { title: 'Price', dataIndex: 'price', key: 'price' },
        { title: 'Total', key: 'total', render: (_, record) => record.quantity * record.price },
    ];

    return (
        <div>
            <h2>Order Details for Order ID: {id}</h2>
            <Table
                dataSource={orderItems}
                columns={columns}
                rowKey="order_item_id"
                loading={loading}
            />
            <div style={{ marginTop: 16 }}>
                <Button
                    type="primary"
                    onClick={openPaymentModal}
                    disabled={isPaid}
                >
                    {isPaid ? 'Order Paid' : 'Pay for Order'}
                </Button>
                <Button
                    type="default"
                    onClick={() => navigate('/customer/dashboard/orders')}
                    style={{ marginLeft: 8 }}
                >
                    Back to Orders
                </Button>
            </div>

            <Modal
                title="Select Payment Method"
                open={paymentModalVisible}
                onCancel={closePaymentModal}
                onOk={handlePayment}
                okText="Pay"
                cancelText="Cancel"
            >
                <Select
                    style={{ width: '100%' }}
                    placeholder="Choose a payment method"
                    onChange={(value) => setSelectedPaymentMethod(value)}
                >
                    <Option value="Credit Card">Credit Card</Option>
                    <Option value="PayPal">PayPal</Option>
                    <Option value="Bank Transfer">Bank Transfer</Option>
                </Select>
            </Modal>
        </div>
    );
}

export default OrderDetails;
