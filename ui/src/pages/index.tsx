
import { Layout, Typography, Button, Card, Row, Col } from 'antd';
//import 'antd/dist/antd.css';

const { Header, Content, Footer } = Layout;
const { Title, Text } = Typography;

const CompanyIntroductionPage = () => {
  return (
    <Layout>
      <Header style={{ background: '#fff', padding: '100px', textAlign: 'center' }}>
        <Title level={2}>斯塔莱特（上海）新材料研究有限公司</Title>
        
      </Header>
      <Content style={{ padding: '5px', background: '#f0f2f5' }}>
        <Row gutter={[16, 16]}>
          <Col span={24}>
            <Card title="About Us" style={{ height: '100%' }}>
              <Text>
              斯塔莱特（上海）新材料研究有限公司是华容县恒兴建材有限公司于2023年5月在上海成立的全资子公司。华容县恒兴建材有限公司是领先的紫外光固化涂料和聚氨酯热熔胶供应商，是国家级专精特新小巨人企业。为加强公司的研发能力，扩展业务领域，在上海成立了斯塔莱特（上海）新材料研究有限公司，从事高性能材料的研究与开发。
              </Text>
            </Card>
          </Col>
          <Col span={12}>
            <Card title="Our Services" style={{ height: '100%' }}>
              <ul>
                <li>Service 1</li>
                <li>Service 2</li>
                <li>Service 3</li>
              </ul>
            </Card>
          </Col>
          <Col span={12}>
            <Card title="Contact Us" style={{ height: '100%' }}>
              <Text>
                Email: info@example.com
                <br />
                Phone: +1 (555) 123-4567
              </Text>
            </Card>
          </Col>
        </Row>
      </Content>
      <Footer style={{ textAlign: 'center', background: '#fff', padding: '20px' }}>
        <Button type="primary">Contact Us</Button>
      </Footer>
    </Layout>
  );
};

export default CompanyIntroductionPage;
