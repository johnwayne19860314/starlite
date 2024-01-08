import React, { useState,useEffect } from 'react';
import { Tabs } from 'antd';

import CollapsibleTable from './components/CollapseTable';
import { listEntryCategory } from '@/services/entry';

interface TabItemProp {
  key: string;
  label: string
}

const CollapsibleTabTable: React.FC = () => {
  
  const [activeTab, setActiveTab] = useState('A');
  const [tabItems, setTabItems] = useState<TabItemProp[]>([]);
  const [load, setLoad] = useState(false)


useEffect(() => {
    // Fetch data from your API endpoint
    const fetchItemsData = async () => {
      try {
        const response = await listEntryCategory({
          IsActive: true,
          
        }); // Replace with your API endpoint
        setTabItems(response.data.items);
        console.log(response)
      } catch (error) {
        console.error('Error fetching data:', error);
      }
    };

    fetchItemsData();
  }, []); // The empty dependency array means this effect will run once after the initial render

  const handleTabChange = (key:any) => {
    console.log("the active tab is ", key)
    setActiveTab(key);
    setLoad(true)
  };
  return (
    <div>
      
      <Tabs
            centered
            activeKey={activeTab}
            onChange={handleTabChange}
            items={tabItems}
          />

      {load  && (
        <CollapsibleTable 
          tabKey={activeTab}
          // loadable = {false}
        />

      )}
    </div>
  );
};

export default CollapsibleTabTable;
