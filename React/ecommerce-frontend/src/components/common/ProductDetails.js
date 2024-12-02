import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import axios from 'axios';
import { Card, Row, Col, Typography, Skeleton, Button, message, List, Rate, Form, Input } from 'antd';
import { useSelector } from 'react-redux';

const { Title, Text, Paragraph } = Typography;

function ProductDetails() {
    const { id } = useParams();
    const [product, setProduct] = useState(null);
    const [reviews, setReviews] = useState([]);
    const [loading, setLoading] = useState(true);
    const [submittingReview, setSubmittingReview] = useState(false);

    const auth = useSelector((state) => state.auth);
    const { user } = auth;

    useEffect(() => {
        axios.get(`/api/products/${id}`)
            .then(res => {
                setProduct(res.data.product);
                setReviews(res.data.reviews || []);
                setLoading(false);
            })
            .catch(err => {
                console.error('Error loading product details:', err);
                setLoading(false);
            });
    }, [id]);

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

    const handleSubmitReview = (values) => {
        setSubmittingReview(true);
        axios.post(`/api/products/${id}/reviews`, values)
            .then(() => {
                message.success('Review submitted successfully');
                // Обновляем список отзывов
                axios.get(`/api/products/${id}`)
                    .then(res => {
                        setReviews(res.data.reviews || []);
                    })
                    .catch(err => {
                        console.error('Error refreshing reviews:', err);
                    });
            })
            .catch(err => {
                console.error('Error submitting review:', err);
                if (err.response && err.response.data) {
                    message.error(err.response.data);
                } else {
                    message.error('Failed to submit review');
                }
            })
            .finally(() => {
                setSubmittingReview(false);
            });
    };

    if (loading) {
        return <Skeleton active />;
    }

    if (!product) {
        return <p>Product not found</p>;
    }

    return (
        <Row justify="center" style={{ padding: '50px 0' }}>
            <Col xs={24} sm={16} md={12} lg={10}>
                <Card
                    hoverable
                    cover={
                        <div
                            style={{
                                height: '300px',
                                display: 'flex',
                                alignItems: 'center',
                                justifyContent: 'center',
                                backgroundColor: '#f5f5f5',
                            }}
                        >
                            <img
                                alt={product.name}
                                src={product.image_url || '/placeholder.jpg'}
                                style={{ maxHeight: '100%', maxWidth: '100%' }}
                            />
                        </div>
                    }
                    actions={[
                        user?.role_name === 'Customer' && (
                            <Button
                                type="primary"
                                onClick={() => handleAddToCart(product.product_id)}
                            >
                                Add to Cart
                            </Button>
                        ),
                    ].filter(Boolean)}
                >
                    <Title level={3}>{product.name}</Title>
                    <Text type="secondary" style={{ display: 'block', marginBottom: '10px' }}>
                        Category: {product.category_name}
                    </Text>
                    <Paragraph style={{ marginTop: '10px' }}>
                        <Text strong>Description:</Text>
                        <br />
                        {product.description || 'No description available.'}
                    </Paragraph>
                    <div style={{ marginTop: '20px' }}>
                        <Text strong style={{ fontSize: '18px' }}>
                            Price: ${product.price}
                        </Text>
                    </div>
                    <div style={{ marginTop: '10px' }}>
                        <Text>Stock: {product.stock}</Text>
                    </div>
                </Card>
            </Col>
            <Col xs={24} sm={16} md={12} lg={10} style={{ marginTop: '20px' }}>
                <Card title="Customer Reviews">
                    {reviews.length > 0 ? (
                        <List
                            dataSource={reviews}
                            renderItem={(review) => (
                                <List.Item>
                                    <List.Item.Meta
                                        title={
                                            <div>
                                                <Rate disabled value={review.rating} />
                                                <Text style={{ marginLeft: '10px' }}>
                                                    {new Date(review.created_at).toLocaleDateString()}
                                                </Text>
                                                <Text style={{ marginLeft: '10px', fontWeight: 'bold' }}>
                                                    {review.author_name}
                                                </Text>
                                            </div>
                                        }
                                        description={review.comment}
                                    />
                                </List.Item>
                            )}
                        />
                    ) : (
                        <Text>No reviews yet.</Text>
                    )}
                </Card>
                {user?.role_name === 'Customer' && (
                    <Card title="Write a Review" style={{ marginTop: '20px' }}>
                        <Form onFinish={handleSubmitReview}>
                            <Form.Item
                                name="rating"
                                rules={[
                                    { required: true, message: 'Please provide a rating' },
                                ]}
                            >
                                <Rate />
                            </Form.Item>
                            <Form.Item
                                name="comment"
                                rules={[
                                    { required: true, message: 'Please write a comment' },
                                    { max: 500, message: 'Comment cannot exceed 500 characters' },
                                ]}
                            >
                                <Input.TextArea rows={4} placeholder="Your review" maxLength={500} />
                            </Form.Item>
                            <Form.Item>
                                <Button type="primary" htmlType="submit" loading={submittingReview}>
                                    Submit Review
                                </Button>
                            </Form.Item>
                        </Form>
                    </Card>
                )}
            </Col>
        </Row>
    );
}

export default ProductDetails;
