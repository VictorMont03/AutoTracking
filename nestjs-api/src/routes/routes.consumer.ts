import { RoutesService } from './routes.service';
import { Controller } from "@nestjs/common";
import { MessagePattern } from "@nestjs/microservices";
import { KafkaContext } from "src/kafka/kafka-context";


@Controller()
export class RouterConsumer {
	constructor(private routesService: RoutesService) { }

	@MessagePattern('freight')
	async updateFreight(payload: KafkaContext) {
		console.log(`Consumer | freight | ${JSON.stringify(payload.messageValue)}`);

		await this.routesService.update(payload.messageValue.route_id, { freight: payload.messageValue.amount });
	}
}