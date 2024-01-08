import  { useEffect, FC } from 'react';
import { Modal, Form, Input, message } from 'antd';

import TextArea from 'antd/es/input/TextArea';
import { newEntryName } from '../entry';

export interface FormValues {
  [name: string]: any;
}
interface EntryCategoryModalProps {
  visible: boolean;
  record: API.DataSourceItem;
  tabKey:string;
  closeHandler: () => void;
  onFinish: (values: FormValues) => void;
  confirmLoading: boolean;
}

const layout = {
  labelCol: { span: 4 },
  wrapperCol: { span: 20 },
};

const DataItemModal: FC<EntryCategoryModalProps> = props => {
  const [form] = Form.useForm();
  const { visible, record, closeHandler, onFinish, confirmLoading,tabKey } = props;

  useEffect(() => {
    if (record === undefined) {
      form.resetFields();
    } else {
      record.code = `${tabKey}000`
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
          <Form.Item label="编码" name="code" initialValue={record.code}>
            <Input />
          </Form.Item>
          <Form.Item label="名称" name="name" initialValue={record.name}>
            <Input />
          </Form.Item>
          <Form.Item label="数量" name="amount" initialValue={record.amount}>
            <Input />
          </Form.Item>
          <Form.Item label="重量" name="weight" initialValue={record.weight}>
            <Input />
          </Form.Item>
          {/* <Form.Item label="note" name="note" initialValue={record.note}>
            <Input />
          </Form.Item> */}
          <Form.Item label="说明" name="note" initialValue={record.note}>
            <TextArea rows={4} />
          </Form.Item>
          
        </Form>
      </Modal>
    </div>
  );
};

export default DataItemModal;
