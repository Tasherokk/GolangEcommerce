import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { Row, Col, Card, Select, Button, message } from 'antd';
import { Link } from 'react-router-dom';
import { useSelector } from 'react-redux'; // Для доступа к состоянию пользователя

const { Option } = Select;

function HomePage() {
    const [products, setProducts] = useState([]);
    const [filteredProducts, setFilteredProducts] = useState([]);
    const [categories, setCategories] = useState([]);
    const [selectedCategory, setSelectedCategory] = useState('all');
    const [sortOrder, setSortOrder] = useState('default');

    const auth = useSelector((state) => state.auth); // Получение состояния пользователя
    const { user } = auth;

    useEffect(() => {
        // Получение продуктов
        axios.get('/api/products')
            .then(res => {
                setProducts(res.data);
                setFilteredProducts(res.data);
            })
            .catch(err => console.error(err));

        // Получение категорий
        axios.get('/api/categories')
            .then(res => setCategories(res.data))
            .catch(err => console.error(err));
    }, []);

    // Фильтрация и сортировка
    useEffect(() => {
        let updatedProducts = [...products];

        // Фильтр по категории
        if (selectedCategory !== 'all') {
            updatedProducts = updatedProducts.filter(
                product => product.category_name === selectedCategory
            );
        }

        // Сортировка по цене
        if (sortOrder === 'asc') {
            updatedProducts.sort((a, b) => a.price - b.price);
        } else if (sortOrder === 'desc') {
            updatedProducts.sort((a, b) => b.price - a.price);
        }

        setFilteredProducts(updatedProducts);
    }, [selectedCategory, sortOrder, products]);

    const handleAddToCart = (productId) => {
        axios.post('/api/cart/items', { product_id: productId, quantity: 1 })
            .then(() => {
                message.success('Product added to cart successfully');
            })
            .catch(err => {
                console.error('Error adding product to cart:', err);
                message.error('Failed to add product to cart');
            });
    };

    return (
        <div>
            <h1>Welcome to Our Shop</h1>

            <div style={{ marginBottom: 16 }}>
                {/* Фильтр по категориям */}
                <Select
                    style={{ width: 200, marginRight: 16 }}
                    value={selectedCategory}
                    onChange={(value) => setSelectedCategory(value)}
                >
                    <Option value="all">All Categories</Option>
                    {categories.map(category => (
                        <Option key={category.category_id} value={category.name}>
                            {category.name}
                        </Option>
                    ))}
                </Select>

                {/* Сортировка по цене */}
                <Select
                    style={{ width: 200 }}
                    value={sortOrder}
                    onChange={(value) => setSortOrder(value)}
                >
                    <Option value="default">Default</Option>
                    <Option value="asc">Price: Low to High</Option>
                    <Option value="desc">Price: High to Low</Option>
                </Select>
            </div>

            <Row gutter={[16, 16]}>
                {filteredProducts.map(product => (
                    <Col key={product.product_id} span={6}>
                        <Card
                            hoverable
                            actions={[
                                user?.role_name === 'Customer' && (
                                    <Button
                                        type="primary"
                                        onClick={() => handleAddToCart(product.product_id)}
                                    >
                                        Add to Cart
                                    </Button>
                                ),
                            ].filter(Boolean)} // Убирает undefined, если пользователь не Customer
                        >
                            <Card.Meta
                                title={<Link to={`/products/${product.product_id}`}>{product.name}</Link>}
                                description={`$${product.price}`}
                            />
                        </Card>
                    </Col>
                ))}
            </Row>
        </div>
    );
}

export default HomePage;
