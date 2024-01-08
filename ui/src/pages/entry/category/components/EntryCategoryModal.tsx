import  { useEffect, FC } from 'react';
import { Modal, Form, Input, message } from 'antd';

import TextArea from 'antd/es/input/TextArea';

export interface FormValues {
  [name: string]: any;
}
interface EntryCategoryModalProps {
  visible: boolean;
  record: API.EntryCategory;
  closeHandler: () => void;
  onFinish: (values: FormValues) => void;
  confirmLoading: boolean;
}

const layout = {
  labelCol: { span: 4 },
  wrapperCol: { span: 20 },
};

const EntrycategoryModal: FC<EntryCategoryModalProps> = props => {
  const [form] = Form.useForm();
  const { visible, record, closeHandler, onFinish, confirmLoading } = props;

  useEffect(() => {
    if (record === undefined) {
      form.resetFields();
    } else {
      
      form.setFieldsValue({
        ...record,
        // create_time: moment(record.create_time),
        // status: Boolean(record.status),
      });
    }
  }, [visible]);

  const onOk = () => {
    form.submit();
  };

  const onFinishFailed = (errorInfo: any) => {
    message.error(errorInfo.errorFields[0].errors[0]);
  };

  return (
    <div>
      <Modal
        title={record ? '更改标签 : ' + record.key : '新增标签'}
        visible={visible}
        onOk={onOk}
        onCancel={closeHandler}
        forceRender
        confirmLoading={confirmLoading}
      >
        <Form
          {...layout}
          name="basic"
          form={form}
          onFinish={onFinish}
          onFinishFailed={onFinishFailed}
          initialValues={{
            status: true,
          }}
        >
          <Form.Item
            label="标签名"
            name="key"
            rules={[{ required: true, message: '请输入标签名!' }]}
          >
            <Input />
          </Form.Item>

          <Form.Item label="说明" name="note">
            <TextArea />
          </Form.Item>
          
        </Form>
      </Modal>
    </div>
  );
};

export default EntrycategoryModal;
