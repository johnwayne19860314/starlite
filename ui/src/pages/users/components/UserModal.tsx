import React, { useEffect, FC } from 'react';
import { Modal, Button, Form, Input, message, DatePicker, Select,Switch } from 'antd';
//import { FormValues } from '../data.d';
//import moment from 'moment';

export interface FormValues {
  [name: string]: any;
}
interface UserModalProps {
  visible: boolean;
  record: API.UserInfo;
  closeHandler: () => void;
  onFinish: (values: FormValues) => void;
  confirmLoading: boolean;
}

const layout = {
  labelCol: { span: 4 },
  wrapperCol: { span: 20 },
};
export const roles = ["undefined","管理员", "专属用户", "内部用户"]
const UserModal: FC<UserModalProps> = props => {
  const [form] = Form.useForm();
  const { visible, record, closeHandler, onFinish, confirmLoading } = props;

  useEffect(() => {
    if (record === undefined) {
      form.resetFields();
    } else {
      record.role = roles[record.role]
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
        title={record ? '更改用户 ID: ' + record.id : '新增用户'}
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
            label="名字"
            name="name"
            rules={[{ required: true, message: '请输入用户名!' }]}
          >
            <Input />
          </Form.Item>

          <Form.Item label="邮箱" name="email">
            <Input />
          </Form.Item>
          {/* <Form.Item label="Create Time" name="create_time">
            <DatePicker showTime />
          </Form.Item> */}
          {/* <Form.Item label="Role" name="role" valuePropName="checked">
            <Switch />
          </Form.Item> */}
          <Form.Item label="角色" name = "role">
            <Select>
              <Select.Option value= {roles[1]}>管理员</Select.Option>
              <Select.Option value={roles[2]}>专属用户</Select.Option>
              <Select.Option value={roles[3]}>内部用户</Select.Option>
            </Select>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default UserModal;
