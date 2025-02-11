import { Client as googleMapsClient, PlaceInputType } from '@googlemaps/google-maps-services-js';
import { Injectable } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';

@Injectable()
export class PlacesService {
	constructor(private readonly googleMapsClient: googleMapsClient, private configService: ConfigService) { }

	async findPlaces(text: string) {
		const { data } = await this.googleMapsClient.findPlaceFromText({
			params: {
				input: text,
				inputtype: PlaceInputType.textQuery,
				fields: ['formatted_address', 'place_id', 'name', "geometry"],
				key: this.configService.get<string>('GOOGLE_MAPS_API_KEY')
			}
		})

		return data
	}
}
