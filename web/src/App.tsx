"use client"

import './App.css'
import MainPage from './components/main/page'
import RootLayout from './layout'

import { useTranslation } from 'react-i18next'


function App() {
  const { t } = useTranslation();
 


  return (
    <>
          <RootLayout>
              <h1 className=" font-normal text-left my-4 text-3xl lg:text-gray-800 ">{t('应用商店')}</h1>
              <div className='w-full'>
                  <MainPage />
              </div>
          
          </RootLayout>
    </>
  )
}


export default App

