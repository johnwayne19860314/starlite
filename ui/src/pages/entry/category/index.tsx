import  { useState, FC,useEffect } from 'react';
import {

  Button,
  Popconfirm,
  message,
} from 'antd';
import ProTable, { ProColumns } from '@ant-design/pro-table';

import EntryCategoryModal , {FormValues} from './components/EntryCategoryModal';

const initEntryCategory = {
  "key":"",
  "label":"",
  "note":""
}

import { addEntryCategory, delEntryCategory, listEntryCategory, updateEntryCategory } from '@/services/entry';

const CategoryListPage: FC = () => {
  const [modalVisible, setModalVisible] = useState(false);
  const [confirmLoading, setConfirmLoading] = useState(false);
  const [record, setRecord] = useState<API.EntryCategory >(initEntryCategory);
  const [dataSource, setDataSource] = useState<API.EntryCategory[]>([]);
  //const [page, setPage] = useState(1)
  const [resetDataSource, setResetDataSource]= useState(false)

  useEffect(() => {
    // Fetch data from your API endpoint
    const fetchItemsData = async () => {
      try {
        const response = await listEntryCategory({
          IsActive: true,
          
        }); // Replace with your API endpoint
        setDataSource(response.data.items);
        console.log(response)
      } catch (error) {
        console.error('Error fetching data:', error);
      }
    };

    fetchItemsData();
  }, [resetDataSource]);

  const columns: ProColumns<API.EntryCategory>[] = [
    // {
    //   title: 'ID',
    //   dataIndex: 'id',
    //   key: 'id',
    // },
    {
      title: '标签',
      dataIndex: 'key',
      key: 'key',
      valueType: 'text',
      //render: (text: any) => {text},
    },
    {
      title: '说明',
      dataIndex: 'note',
      key: 'note',
    },
    
    {
      title: '操作',
      key: 'action',
      valueType: 'option',
      render: (text: any, record: API.EntryCategory) => [
        <a
          onClick={() => {
            editHandler(record);
          }}
        >
          更改
        </a>,
        <Popconfirm
          title="确定要删除此标签吗？与之相关的原料清单都会被删除"
          onConfirm={() => {
            deleteHandler(record.key);
          }}
          okText="确认"
          cancelText="取消"
        >
          <a>删除</a>
        </Popconfirm>,
      ],
    },
  ];

  const editHandler = (record: API.EntryCategory) => {
    setModalVisible(true);
    setRecord(record);
  };

  const deleteHandler = (key: string) => {
    const delUserRecord = async () => {
      try {
        const response = await delEntryCategory(key); // Replace with your API endpoint
       
        console.log(response)
        resetHandler()
      } catch (error) {
        console.error('Error fetching data:', error);
      }
    };

    delUserRecord();
  };

  const closeHandler = () => {
    setModalVisible(false);
  };

  const onFinish = async (values: FormValues) => {
    //const values = await form.validateFields();
    console.log("=======", values)
    setConfirmLoading(true);
    let category = "";

    if (record) {
      category = record.key;
    }

    let serviceFun;
    if (category) {
      serviceFun = updateEntryCategory;
    } else {
      serviceFun = addEntryCategory;
    }
    
    const result = await serviceFun({
      entryCategory:{
        category:values.key,
        note:values.note,
      }
    });
    if (result) {
      setModalVisible(false);
      message.success(`${category === "" ? '新增' : '更改'} 成功。`);
      resetHandler();
    } else {
      message.error(`${category === "" ? '新增' : '更改'} 失败。`);
    }
    setConfirmLoading(false);
  };

  const addHandler = () => {
    setModalVisible(true);
    setConfirmLoading(false)
    setRecord(initEntryCategory);
  };

  const resetHandler = () => {
    if (resetDataSource) {
      setResetDataSource(false)
    }else{
      setResetDataSource(true)
    }
    
  };


  return (
    <div className="list-table">
      <ProTable
        columns={columns}
        dataSource={dataSource}
        rowKey="key"
        search={false}
        
        options={{
          density: true,
          fullScreen: true,
          reload: () => {
            resetHandler();
          },
          setting: true,
        }}
        headerTitle="标签类型列表"
        toolBarRender={() => [
          <Button type="primary" onClick={addHandler}>
            新增
          </Button>,
          <Button onClick={resetHandler}>刷新</Button>,
        ]}
      />
      
      <EntryCategoryModal
        visible={modalVisible}
        closeHandler={closeHandler}
        record={record}
        onFinish={onFinish}
        confirmLoading={confirmLoading}
      ></EntryCategoryModal>
    </div>
  );
};


export default CategoryListPage;
