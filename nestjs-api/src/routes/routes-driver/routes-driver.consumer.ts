import { Controller } from "@nestjs/common";
import { MessagePattern } from "@nestjs/microservices";
import { KafkaContext } from "src/kafka/kafka-context";
import { RoutesDriverService } from "./routes-driver.service";
import { HttpService } from "@nestjs/axios";


@Controller()
export class RouterDriverConsumer {
	constructor(private httpService: HttpService) { }

	@MessagePattern('simulation')
	async driverMoved(payload: KafkaContext) {
		console.log(`Consumer | simulation | ${JSON.stringify(payload.messageValue)}`);

		await this.httpService.axiosRef.post(`http://localhost:${process.env.PORT}/routes/${payload.messageValue.route_id}/process-route`, {
			lat: payload.messageValue.lat,
			lng: payload.messageValue.lng
		});
	}
}