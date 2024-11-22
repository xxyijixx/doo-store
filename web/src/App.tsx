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
              <h1 className=" font-bold text-left my-4 text-5xl md:text-5xl lg:text-3xl">{t('云应用')}</h1>
              <div className='w-full'>
                  <MainPage />
              </div>
          
          </RootLayout>
    </>
  )
}


export default App

