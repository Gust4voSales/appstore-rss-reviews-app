import { StatusBar } from 'expo-status-bar';

import './global.css';
import { Container } from 'components/Container';
import { Text } from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';

export default function App() {
  return (
    <SafeAreaView className="flex-1">
      <Container>
        <Text className="my-auto text-center text-2xl font-bold text-black">Hello</Text>
      </Container>
      <StatusBar style="auto" />
    </SafeAreaView>
  );
}
