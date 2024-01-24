import  { useEffect, FC } from 'react';
import { Modal, Form, Input, message } from 'antd';

import TextArea from 'antd/es/input/TextArea';


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
      if (visible){
        if (record.name == "") {
          record.code = getCode(record.code)
        }
        
        form.setFieldsValue({
          ...record,
          // create_time: moment(record.create_time),
          // status: Boolean(record.status),
        });
      }
      
    }
  }, [visible]);

  const getCode = (code:string) :string =>  {
    
    // let top = "X1004"
    // let codeType = code.slice(0,1)
    console.log("come to the getcode function")
    if (tabKey == "X") {
      if (code == "") {
        return "X1001"
      }
      let base = 7
      let topCode = parseInt(code.slice(1,code.length))
      let head = Math.floor(topCode/1000)
      let tail = topCode%1000
      if (head == base) {
        head = 1
        tail += 3
  
      }else{
        head +=1
      }
      return `X${head*1000+tail}` 
    } else if (tabKey == "X8"){
      if (code == "") {
        return "X8002"
      }
      
    }else if (tabKey == "S"){
      if (code == "") {
        return "S001"
      }
      
    } else if (tabKey == "C") {
      if (code == "") {
        return "C002"
      }
     
    }
    let topCode = parseInt(code.slice(tabKey.length,code.length))
    let tmp = (1000+topCode+2)
    return (""+tmp).replace("1",tabKey)
    
  }

  const onOk = () => {
    form.submit();
  };

  const onFinishFailed = (errorInfo: any) => {
    message.error(errorInfo.errorFields[0].errors[0]);
  };

  return (
    <div>
      <Modal
        title={record.name ? '更改标签 : ' + record.code : '新增标签'}
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
          <Form.Item label="编码" name="code" initialValue={record.code} required>
            <Input />
          </Form.Item>
          <Form.Item label="名称" name="name" initialValue={record.name} required>
            <Input />
          </Form.Item>
          {/* <Form.Item label="数量" name="amount" initialValue={record.amount}>
            <Input />
          </Form.Item> */}
          <Form.Item label="重量" name="weight" initialValue={record.weight} required>
            <Input />
          </Form.Item>
          <Form.Item label="供应商名称" name="supplier" initialValue={record.supplier}>
            <Input />
          </Form.Item>
          <Form.Item label="供应商联系方式" name="supplierContactInfo" initialValue={record.supplierContactInfo}>
            <Input />
          </Form.Item>
          {/* <Form.Item label="note" name="note" initialValue={record.note}>
            <Input />
          </Form.Item> */}
          <Form.Item label="说明" name="note" initialValue={record.note} required>
            <TextArea rows={4} />
          </Form.Item>
          
        </Form>
      </Modal>
    </div>
  );
};

export default DataItemModal;
