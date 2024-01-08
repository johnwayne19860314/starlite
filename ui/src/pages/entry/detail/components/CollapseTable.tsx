import React, { useState,useEffect } from 'react';
import { Table, Collapse, Button, Popconfirm, Form, Input,Space,message } from 'antd';
import { addEntry, listEntries, updateEntry,delEntry } from '@/services/entry';
import DataItemModal, {FormValues} from './DataItemModal';
import { newEntryName,initDataSourceItem } from '../entry.d';

const { Panel } = Collapse;
const { TextArea } = Input;



interface CollapseTableProp {
  tabKey: string;
  //loadable: boolean;
}


const CollapsibleTable: React.FC<CollapseTableProp> = props => {
  const [modalVisible, setModalVisible] = useState(false);
  const [confirmLoading, setConfirmLoading] = useState(false);
  const [record, setRecord] = useState<API.DataSourceItem >(initDataSourceItem);
  const [resetDataSource, setResetDataSource]= useState(false)
  
  const [expandedRowKey, setExpandedRowKey] = useState<string | undefined>(undefined);
  const [dataSource, setDataSource] = useState<API.DataSourceItem[]>([]);
  const [editingKey, setEditingKey] = useState<string | undefined>(undefined);
  const [filters, setFilters] = useState<any>({});
  
  const { tabKey } = props;
  console.log("=====tabkey is ======",props)

  useEffect(() => {
    // Fetch data from your API endpoint
    const fetchData = async () => {
      try {
        const response = await listEntries({
          IsActive: true,
          Category:tabKey,
          Offset: 0,
          Limit: 9999999,
        }); // Replace with your API endpoint

        setDataSource(response.data.entries||[]);
        console.log(response)
      } catch (error) {
        console.error('Error fetching data:', error);
      }
    };

    fetchData();
  }, [tabKey,resetDataSource]); // The empty dependency array means this effect will run once after the initial render

  const closeHandler = () => {
    setModalVisible(false);
  };

  const onFinish = async (values: FormValues) => {
    //const values = await form.validateFields();
    console.log("=======", values)
    setConfirmLoading(true);
    
    const entry = {
      name : values.name,
      code : values.code,
      codeCategory : tabKey,
      amount :  parseInt(values.amount),
      weight : parseFloat(values.weight),
      note : values.note,
      isActive : true
    }
    const result = await addEntry({entry});
    if (result) {
      setModalVisible(false);
      message.success(`新增成功。`);
      resetHandler();
    } else {
      message.error(`新增失败。`);
    }
    setConfirmLoading(false);
  };

  // const addHandler = () => {
  //   setModalVisible(true);
  //   setConfirmLoading(false)
  //   setRecord(initEntryCategory);
  // };

  const resetHandler = () => {
    if (resetDataSource) {
      setResetDataSource(false)
    }else{
      setResetDataSource(true)
    }
    
  };

  const handleExpand = (key: string) => {
    setExpandedRowKey((prevKey) => (prevKey === key ? undefined : key));
  };

  const handleEdit = (record: API.DataSourceItem) => {
    setEditingKey(record.key);
    setExpandedRowKey(record.key);
    //handleExpand(record.key)
  };

  const handleDelete = (code: string) => {
    const delEntryFunc = async () => {
      try {
       await delEntry(code); // Replace with your API endpoint
        
        setDataSource((prevDataSource) => prevDataSource.filter((item) => item.code !== code));
        setEditingKey(undefined);
        setExpandedRowKey(undefined);
      } catch (error) {
        console.error('Error fetching data:', error);
      }
    };
    delEntryFunc();

    
  };

  // const handleSave = (form: any, key: string) => {
  //   form.validateFields((error: any, values: any) => {
  //     if (error) {
  //       return;
  //     }
      
  //   });
  // };

  const handleSave = async (form: any, key: string) => {
    try {
      const values = await form.validateFields();
      const entry = {
        name : values.name,
        code : values.code,
        codeCategory : tabKey,
        amount :  parseInt(values.amount),
        weight : parseFloat(values.weight),
        note : values.note,
        isActive : true
      }
      const saveEntry = async () => {
        try {
          if (entry.name == newEntryName) {
            await addEntry({
              entry
            }); 
          }else {
            await updateEntry(entry)
          }
          let newData = [...dataSource];
          let index = newData.findIndex((item) => key === item.key);
          if (index > -1) {
            newData[index] = { ...newData[index], ...values };
            setDataSource(newData);
            setEditingKey(undefined);
            setExpandedRowKey(undefined);
          }
        } catch (error) {
          console.error('Error fetching data:', error);
        }
      };
  
      saveEntry();
      
      // Your asynchronous logic here
    } catch (error) {
      console.error('Validation failed:', error);
    }
  };

  const columns = [
    
    {
      title: '编号',
      dataIndex: 'code',
      key: 'code',
    },
    // {
    //   title: 'Name',
    //   dataIndex: 'name',
    //   key: 'name',
    // },
    // {
    //   title: 'Amount',
    //   dataIndex: 'amount',
    //   key: 'amount',
    // },
    // {
    //   title: 'Weight',
    //   dataIndex: 'weight',
    //   key: 'weight',
    // },
    // {
    //   title: 'Note',
    //   dataIndex: 'note',
    //   key: 'note',
    // },
    {
      title: '操作',
      key: 'actions',
      render: (record:any) => (
        <span>
          <Button type="primary" style={{ marginRight: 8 }} onClick={() => handleEdit(record)}>
            更改
          </Button>
          <Popconfirm
            title="Are you sure you want to delete this row?"
            onConfirm={() => handleDelete(record.code)}
          >
            <Button danger>删除</Button>
          </Popconfirm>
        </span>
      ),
    },
  ];
  const addRow = () => {
    
    const newKey = (dataSource.length + 1).toString();
    const newRow: API.DataSourceItem = {
      key: newKey,
      id :0,
      name: newEntryName,
      code: `${tabKey}${newKey}`,
      amount: 0,
      weight:0.0,
      note: '新的描述',
    };
    setDataSource([newRow, ...dataSource]);
    setEditingKey(newKey);
    setExpandedRowKey(undefined);
  };

  const addHandler = () => {
    setModalVisible(true);
    setConfirmLoading(false)
    setRecord(initDataSourceItem);
  };
  const handleFilterChange = (columnKey: string, value: string) => {
    setFilters((prevFilters:any) => ({ ...prevFilters, [columnKey]: value }));
  };
  const filteredDataSource = dataSource.filter((item) => {
    return Object.keys(filters).every((key) => {
      if (key === 'amount') {
        return item[key] >= parseInt(filters[key], 10);
      }
      if (key === 'weight') {
        return item[key] >= parseFloat(filters[key]);
      }
      //return item[key].toLowerCase().includes(filters[key].toLowerCase());
      return String(item[key]).toLowerCase().includes(filters[key].toLowerCase());
    });
  });
  const isEditing = (record: API.DataSourceItem) => record.key === editingKey;

  const editColumns = columns.map((col) => {
    // if (!col.editable) {
    //   return col;
    // }
    const editableCol = col as {
      title: string;
      dataIndex: string;
      key: string;
      editable?: boolean; // Add editable property
      onCell: (record: API.DataSourceItem) => {
        record: API.DataSourceItem;
        inputType: 'text';
        dataIndex: string;
        title: string;
        editing: boolean;
      };
    };
  
    return {
      ...editableCol,
      onCell: (record: API.DataSourceItem) => ({
        record,
        inputType: 'text',
        dataIndex: editableCol.dataIndex,
        title: editableCol.title,
        editing: isEditing(record),
      }),
    };

  });
  const PanelContent: React.FC<{ record: API.DataSourceItem }> = ({ record }) => (
    <div>
      {isEditing(record) ? (
        <Form form={form} onFinish={() => handleSave(form, record.key)}>
          {/* <Form.Item label="编码" name="code"  initialValue={record.code}>
            <Input />
          </Form.Item> */}
          <Form.Item label="名称" name="name" initialValue={record.name}>
            <Input />
          </Form.Item>
          {/* <Form.Item label="数量" name="amount" initialValue={record.amount}>
            <Input />
          </Form.Item> */}
          <Form.Item label="重量" name="weight" initialValue={record.weight}>
            <Input />
          </Form.Item>
          {/* <Form.Item label="note" name="note" initialValue={record.note}>
            <Input />
          </Form.Item> */}
          <Form.Item label="说明" name="note" initialValue={record.note}>
            <TextArea rows={4} />
          </Form.Item>
          <Button type="primary" htmlType="submit">
            保存
          </Button>
        </Form>
      ) : (
        <>
        {/* <p>details go here: </p> */}
        <p>名称 : {record.name} </p>
        <p>数量 : {record.amount}</p>
        <p>重量 : {record.weight}</p>
        <div>说明 : {record.note}</div>
        </>
        
      )}
    </div>
  );

  const [form] = Form.useForm();

  return (
    <div>
      <Space style={{ marginBottom: 16 }}>
          <Button type="primary" onClick={addHandler}>
            新增原料
          </Button>
          <Input
            placeholder="用原料名称过滤"
            onChange={(e) => handleFilterChange('name', e.target.value)}
            value={filters.name || ''}
          />
          <Input
            placeholder="用原料编码过滤"
            onChange={(e) => handleFilterChange('code', e.target.value)}
            value={filters.age || ''}
          />
        </Space>
        <Table
          dataSource={filteredDataSource}
          columns={[
            ...editColumns,
            {
              title: '详情',
              key: 'details',
              render: (_, record) => (
                <Collapse
                  activeKey={expandedRowKey === record.key ? [record.key] : []}
                  onChange={() => handleExpand(record.key)}
                >
                  <Panel key={record.key} header={`展开`}>
                    <PanelContent record={record} />
                  </Panel>
                </Collapse>
              ),
            },
          ]}
          rowKey="key"
          pagination={{
            pageSize: 5,
            hideOnSinglePage: true,
            total:dataSource.length
          }}
        />

        <DataItemModal
          visible={modalVisible}
          closeHandler={closeHandler}
          record={record}
          tabKey={tabKey}
          onFinish={onFinish}
          confirmLoading={confirmLoading}
        ></DataItemModal>
      
    </div>
  );
};

export default CollapsibleTable;
